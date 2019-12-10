package counter

type Counter interface {
	CounterIncrement(key string, amount int64) (int64, error)
	CounterGet(key string) (int64, error)
}
