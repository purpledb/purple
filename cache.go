package strato

const defaultTtl = 5

type (
	Cache interface {
		CacheGet(key string) (string, error)
		CacheSet(key, value string, ttl int) error
	}

	CacheItem struct {
		Value      string
		Timestamp  int64
		TTLSeconds int
	}
)
