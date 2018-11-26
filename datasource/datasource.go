package datasource

type DataSource interface {
	Value(key string) (interface{}, error)
}
