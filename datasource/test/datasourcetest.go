package datasourcetest

import (
	"github.com/kdimak/go_datasource/datasource"
	"testing"
)

type TestConfig struct {
	DS datasource.WritableDataSource
}

func TestWritableDataSource(t *testing.T, tc TestConfig) {
	tests := []struct {
		title string
		run   func(t *testing.T, c TestConfig)
	}{
		{"should store and read the key", testStoreAndGet},
		{"should return nil on not existent key", testGetNotExistent},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.run(t, tc)
		})
	}
}

func testStoreAndGet(t *testing.T, tc TestConfig) {
	err := tc.DS.Store("foo", "bar")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	actual, err := tc.DS.Value("foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if actual != "bar" {
		t.Errorf("reading existent key from cache. expected: %v, actual: %v", "bar", actual)
	}
}

func testGetNotExistent(t *testing.T, tc TestConfig) {
	actual, err := tc.DS.Value("unknown")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if actual != nil {
		t.Errorf("reading not existent key from cache. expected: %v, actual: %v", nil, actual)
	}
}
