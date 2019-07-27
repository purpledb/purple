package strato

type Counter interface {
	IncrementCounter(key string, amount int32)
	GetCounter(key string) int32
}
