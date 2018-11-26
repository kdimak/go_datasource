package datasourcetest

import (
	"github.com/kdimak/go_datasource/datasource"
	"testing"
)

func TestDatabase_Value(t *testing.T) {
	TestWritableDataSource(t, TestConfig{DS: datasource.NewEmptyDistributedCache()})
	TestWritableDataSource(t, TestConfig{DS: datasource.NewEmptyDatabase()})
}
