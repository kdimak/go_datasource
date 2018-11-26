package datasourcetest

import (
	"github.com/kdimak/go_datasource/datasource"
	"reflect"
	"testing"
)

func TestRemoveLeastRecentlyAddedElement(t *testing.T) {
	lruCache := datasource.NewLruCache(2)

	store(lruCache, "key1", "value1", false, t)
	store(lruCache, "key2", "value2", false, t)
	store(lruCache, "key3", "value3", true, t)

	checkGet(lruCache, "key2", "value2", true, t)
	checkGet(lruCache, "key3", "value3", true, t)
	checkGet(lruCache, "key1", nil, false, t)
}

func TestRemoveLeastRecentlyAccessedElement(t *testing.T) {
	lruCache := datasource.NewLruCache(2)

	store(lruCache, "key1", "value1", false, t)
	store(lruCache, "key2", "value2", false, t)
	checkGet(lruCache, "key1", "value1", true, t)
	store(lruCache, "key3", "value3", true, t)

	checkGet(lruCache, "key1", "value1", true, t)
	checkGet(lruCache, "key3", "value3", true, t)
	checkGet(lruCache, "key2", nil, false, t)
}

func store(cache datasource.LruCache, key, value string, expectedEvicted bool, t *testing.T) {
	if actualEvicted := cache.Store(key, value); actualEvicted != expectedEvicted {
		t.Errorf("storing in LRU cache. expected evicted: %v, actual: %v", expectedEvicted, actualEvicted)
	}
}

func checkGet(cache datasource.LruCache, key string, expectedValue interface{}, expectedPresent bool, t *testing.T) {
	if actualValue, actualPresent := cache.Value(key); actualPresent != expectedPresent {
		t.Errorf("get value from LRU csache: expected present: %v, actual: %v", expectedPresent, actualPresent)
	} else {
		if !reflect.DeepEqual(expectedValue, actualValue) {
			t.Errorf("reading cache. expected: %v, actual: %v", expectedValue, actualValue)
		}
	}
}
