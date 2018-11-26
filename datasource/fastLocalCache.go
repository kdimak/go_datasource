package datasource

import (
	"sync"
)

const fastLocalCacheLruSize = 100
const remoteDataSourcesCount = 2

type FastLocalCache struct {
	localCache       LruCache
	distCache        WritableDataSource
	readonlyDatabase DataSource
	mutex            sync.Mutex
}

type CacheItem struct {
	receive chan CacheResult
	done    chan bool
	value   interface{}
	err     error
}

type CacheResult struct {
	value   interface{}
	err     error
	onApply func()
}

func NewFastLocalDataSource(database DataSource, cache WritableDataSource) DataSource {
	return &FastLocalCache{
		localCache:       NewSyncLruCache(fastLocalCacheLruSize),
		distCache:        cache,
		readonlyDatabase: database,
	}
}

func (flc *FastLocalCache) Value(key string) (interface{}, error) {
	flc.mutex.Lock()

	var item *CacheItem
	value, present := flc.readFromLocalCache(key)
	if present {
		item = value.(*CacheItem)
	} else {
		item = &CacheItem{
			receive: make(chan CacheResult),
			done:    make(chan bool),
		}
		flc.storeToLocalCache(key, item)

		go func(item *CacheItem) {
			received := 0
			for received < remoteDataSourcesCount {
				// todo it's worth to apply select on timeout to protect from hanging datasource requests
				result := <-item.receive
				received++
				item.value = result.value
				item.err = result.err
				if item.value != nil && item.err == nil {
					if result.onApply != nil {
						result.onApply()
					}
					break
				}
			}
			close(item.done)
		}(item)

		go func(item *CacheItem) {
			value, err := flc.readFromRemoteCache(key)
			item.receive <- CacheResult{value: value, err: err}
		}(item)

		go func(item *CacheItem) {
			value, err := flc.readFromDatabase(key)
			item.receive <- CacheResult{value: value, err: err, onApply: func() {
				flc.storeToDistCache(key, value)
			}}
		}(item)
	}
	flc.mutex.Unlock()

	<-item.done

	return item.value, item.err
}

func (flc *FastLocalCache) readFromLocalCache(key string) (interface{}, bool) {
	return flc.localCache.Value(key)
}

func (flc *FastLocalCache) readFromRemoteCache(key string) (interface{}, error) {
	return flc.distCache.Value(key)
}

func (flc *FastLocalCache) readFromDatabase(key string) (interface{}, error) {
	return flc.readonlyDatabase.Value(key)
}

func (flc *FastLocalCache) storeToDistCache(key string, value interface{}) {
	go func() {
		flc.distCache.Store(key, value)
	}()
}

func (flc *FastLocalCache) storeToLocalCache(key string, value interface{}) {
	flc.localCache.Store(key, value)
}
