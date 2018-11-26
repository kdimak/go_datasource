package datasourcetest

import (
	"github.com/kdimak/go_datasource/datasource"
	"testing"
)

type MockWritableDataSource struct {
	valueCallsCount int
	storeCallsCount int
	returnValue     interface{}
}

func (mds *MockWritableDataSource) Value(key string) (interface{}, error) {
	mds.valueCallsCount++
	return mds.returnValue, nil
}

func (mds *MockWritableDataSource) Store(key string, value interface{}) error {
	mds.storeCallsCount++
	return nil
}

func TestUseDatabaseAndCacheOnFirstKeyAccess(t *testing.T) {
	mDatabase := MockWritableDataSource{}
	mCache := MockWritableDataSource{}
	localDataSource := datasource.NewLocalDataSource(&mDatabase, &mCache)

	if _, err := localDataSource.Value("key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if mCache.valueCallsCount != 1 {
		t.Errorf("reading local cache. expected request to distributed cache")
	}
	if mDatabase.valueCallsCount != 1 {
		t.Errorf("reading local cache. expected request to database")
	}
}

func TestGetValueFromCache(t *testing.T) {
	expectedValue := "cache value"
	mDatabase := MockWritableDataSource{}
	mCache := MockWritableDataSource{returnValue: expectedValue}
	localDataSource := datasource.NewLocalDataSource(&mDatabase, &mCache)

	if actualValue, err := localDataSource.Value("key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if actualValue != expectedValue {
		t.Errorf("reading cache. expected: %v, actual: %v", expectedValue, actualValue)
	}
}

func TestGetValueFromDatabase(t *testing.T) {
	expectedValue := "db value"
	mDatabase := MockWritableDataSource{returnValue: expectedValue}
	mCache := MockWritableDataSource{}
	localDataSource := datasource.NewLocalDataSource(&mDatabase, &mCache)

	if actualValue, err := localDataSource.Value("key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if actualValue != expectedValue {
		t.Errorf("reading cache. expected: %v, actual: %v", expectedValue, actualValue)
	}
}

func TestCachesValueOnSubsequentKeyAccesses(t *testing.T) {
	mDatabase := MockWritableDataSource{returnValue: "mock value"}
	mCache := MockWritableDataSource{returnValue: "mock value"}
	localDataSource := datasource.NewLocalDataSource(&mDatabase, &mCache)

	if _, err := localDataSource.Value("key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mCacheValueCalls := mCache.valueCallsCount
	mDatabaseValueCalls := mDatabase.valueCallsCount
	if _, err := localDataSource.Value("key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if mCache.valueCallsCount != mCacheValueCalls {
		t.Errorf("reading local cache. expected no extra request to distributed cache")
	}
	if mDatabase.valueCallsCount != mDatabaseValueCalls {
		t.Errorf("reading local cache. expected no extra request to database")
	}
}

func TestCachesNotExistentKeyAccess(t *testing.T) {
	mDatabase := MockWritableDataSource{}
	mCache := MockWritableDataSource{}
	localDataSource := datasource.NewLocalDataSource(&mDatabase, &mCache)

	if value, err := localDataSource.Value("not-existent-key"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		if value != nil {
			t.Errorf("reading local cache. expected nil value")
		}
		if mCache.valueCallsCount != 1 {
			t.Errorf("reading local cache. expected request to distributed cache")
		}
		if mDatabase.valueCallsCount != 1 {
			t.Errorf("reading local cache. expected request to database")
		}
	}

	if value, err := localDataSource.Value("not-existent-key"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		if value != nil {
			t.Errorf("reading local cache. expected nil value")
		}
		if mCache.valueCallsCount != 1 {
			t.Errorf("reading local cache. expected no extra request to distributed cache")
		}
		if mDatabase.valueCallsCount != 1 {
			t.Errorf("reading local cache. expected no extra request to database")
		}
	}

}
