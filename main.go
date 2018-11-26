package main

import (
	"fmt"
	"github.com/kdimak/go_datasource/datasource"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
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

func main() {
	defer elapsedTime(time.Now(), "All")

	cache := datasource.NewFastLocalDataSource(
		datasource.NewPopulatedDatabase(content),
		datasource.NewEmptyDistributedCache())

	var requestsCntr int32
	var wg sync.WaitGroup
	wg.Add(10 * 50)

	requestNo := func() int32 {
		return atomic.AddInt32(&requestsCntr, 1)
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := 0; j < 50; j++ {
				start := time.Now()
				requestKey := randomKey()
				//requestKey := keys[i]
				requestValue, _ := cache.Value(requestKey)
				fmt.Printf("[%d] Request '%s', response '%v', time: %v\n",
					requestNo(), requestKey, requestValue, time.Since(start))
				wg.Done()
			}
		}(i)
	}
	wg.Wait()
}

func randomKey() string {
	return keys[rand.Int()%entriesCount]
}

func elapsedTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
