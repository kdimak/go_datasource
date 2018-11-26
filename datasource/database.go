package datasource

import (
	"sync"
	"time"
)

type Database struct {
	data sync.Map
}

func (db *Database) Value(key string) (interface{}, error) {
	// simulate 500ms roundtrip to the distributed cache
	time.Sleep(500 * time.Millisecond)

	value, _ := db.data.Load(key)
	return value, nil
}

func (db *Database) Store(key string, value interface{}) error {
	// simulate 500ms roundtrip to the distributed cache
	time.Sleep(500 * time.Millisecond)

	db.data.Store(key, value)
	return nil
}

func (db *Database) fillWith(content map[string]interface{}) {
	for k, v := range content {
		db.data.Store(k, v)
	}
	return
}

func NewEmptyDatabase() *Database {
	return &Database{}
}

func NewPopulatedDatabase(content map[string]interface{}) *Database {
	db := Database{}
	db.fillWith(content)
	return &db
}
