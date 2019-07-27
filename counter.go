package strato

type Counter interface {
	Increment(key string, amount int)
}
