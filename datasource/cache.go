package datasource

import (
	"sync"
	"time"
)

type DistributedCache struct {
	data sync.Map
}

func (dc *DistributedCache) Value(key string) (interface{}, error) {
	// simulate 100ms roundtrip to the distributed cache
	time.Sleep(100 * time.Millisecond)

	value, _ := dc.data.Load(key)
	return value, nil
}

func (dc *DistributedCache) Store(key string, value interface{}) error {
	// simulate 100ms roundtrip to the distributed cache
	time.Sleep(100 * time.Millisecond)

	dc.data.Store(key, value)
	return nil
}

func NewEmptyDistributedCache() *DistributedCache {
	return &DistributedCache{}
}
