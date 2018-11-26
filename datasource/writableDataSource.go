package datasource

type WritableDataSource interface {
	DataSource
	Store(key string, value interface{}) error
}
