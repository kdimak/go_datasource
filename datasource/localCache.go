package datasource

const localCacheLruSize = 100

type LocalCache struct {
	localCache       LruCache
	distCache        WritableDataSource
	readonlyDatabase DataSource
}

func NewLocalDataSource(database DataSource, cache WritableDataSource) DataSource {
	return &LocalCache{
		localCache:       NewSyncLruCache(localCacheLruSize),
		distCache:        cache,
		readonlyDatabase: database,
	}
}

type ValueResult struct {
	value interface{}
	err   error
}

func (lc *LocalCache) Value(key string) (interface{}, error) {
	var localCacheErr error

	// Check local cache
	if value, present := lc.readFromLocalCache(key); present {
		return value, nil
	}

	// Check distributed cache
	if value, err := lc.readFromRemoteCache(key); err == nil && value != nil {
		lc.storeToLocalCache(key, value)
		return value, nil
	}

	// Check database
	if value, err := lc.readFromDatabase(key); err == nil && value != nil {
		lc.storeToLocalCache(key, value)
		lc.storeToDistCache(key, value)
		return value, nil
	}

	lc.storeToLocalCache(key, nil)
	return nil, localCacheErr
}

func (lc *LocalCache) readFromLocalCache(key string) (interface{}, bool) {
	return lc.localCache.Value(key)
}

func (lc *LocalCache) readFromRemoteCache(key string) (interface{}, error) {
	return lc.distCache.Value(key)
}

func (lc *LocalCache) readFromDatabase(key string) (interface{}, error) {
	return lc.readonlyDatabase.Value(key)
}

func (lc *LocalCache) storeToDistCache(key string, value interface{}) {
	go func() {
		lc.distCache.Store(key, value)
	}()
}

func (lc *LocalCache) storeToLocalCache(key string, value interface{}) {
	lc.localCache.Store(key, value)
}
