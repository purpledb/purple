package strato

type Counter interface {
	CounterIncrement(key string, amount int32)
	CounterGet(key string) int32
}
