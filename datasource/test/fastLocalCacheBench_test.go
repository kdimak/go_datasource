package datasourcetest

import (
	"github.com/kdimak/go_datasource/datasource"
	"math/rand"
	"sync"
	"testing"
)

var (
	content = map[string]interface{}{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
		"key6": "value6",
		"key7": "value7",
		"key8": "value8",
		"key9": "value9",
	}
	keys         = []string{"key0", "key1", "key2", "key3", "key4", "key5", "key6", "key7", "key8", "key9"}
	entriesCount = len(keys)
)

var value interface{}

func BenchmarkConcurrentAccess(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		cache := datasource.NewFastLocalDataSource(
			datasource.NewPopulatedDatabase(content),
			datasource.NewEmptyDistributedCache())

		var wg sync.WaitGroup
		wg.Add(10 * 50)
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 50; j++ {
					requestKey := randomKey()
					b.StartTimer()
					value, _ = cache.Value(requestKey)
					b.StopTimer()
					wg.Done()
				}
			}()
		}
		wg.Wait()
	}

}
func randomKey() string {
	return keys[rand.Int()%entriesCount]
}
