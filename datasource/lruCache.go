package datasource

import (
	"container/list"
	"sync"
)

type LruCache interface {
	Value(key string) (value interface{}, present bool)
	Store(key string, value interface{}) (evicted bool)
}

type SimpleLruCache struct {
	size      int
	usageList *list.List
	content   map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func NewLruCache(size int) LruCache {
	return &SimpleLruCache{
		size:      size,
		usageList: list.New(),
		content:   make(map[string]*list.Element),
	}
}

func (lc *SimpleLruCache) Value(key string) (value interface{}, present bool) {
	if e, ok := lc.content[key]; ok {
		lc.usageList.MoveToFront(e)
		value = e.Value.(*entry).value
		return value, true
	}
	return nil, false
}

func (lc *SimpleLruCache) Store(key string, value interface{}) (evicted bool) {
	// Check if such key is already stored.
	if e, ok := lc.content[key]; ok {
		lc.usageList.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}

	// Store the new key.
	e := &entry{key, value}
	lc.content[key] = lc.usageList.PushFront(e)

	// Cut the cache if needed.
	if lc.usageList.Len() > lc.size {
		toRemove := lc.usageList.Back()
		lc.usageList.Remove(toRemove)
		delete(lc.content, toRemove.Value.(*entry).key)
		evicted = true
	}

	return
}

type SyncLruCache struct {
	mutex    sync.RWMutex
	lruCache LruCache
}

func (slc *SyncLruCache) Value(key string) (interface{}, bool) {
	slc.mutex.Lock()
	defer slc.mutex.Unlock()
	return slc.lruCache.Value(key)
}

func (slc *SyncLruCache) Store(key string, value interface{}) bool {
	slc.mutex.Lock()
	defer slc.mutex.Unlock()
	return slc.lruCache.Store(key, value)
}

func NewSyncLruCache(size int) LruCache {
	return &SyncLruCache{lruCache: NewLruCache(size)}
}
