package strato

type Counter interface {
	CounterIncrement(key string, amount int32) error
	CounterGet(key string) (int32, error)
}
