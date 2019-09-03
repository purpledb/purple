package cache

const DefaultTtl = 5

type (
	Cache interface {
		CacheGet(key string) (string, error)
		CacheSet(key, value string, ttl int32) error
	}

	Item struct {
		Value      string
		Timestamp  int64
		TTLSeconds int32
	}
)
