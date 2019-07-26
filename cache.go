package strato

import "strato/proto"

const defaultTtl = 5

type (
	Cache interface {
		CacheGet(key string) (string, error)
		CacheSet(key, value string, ttl int32) error
	}

	CacheItem struct {
		Value      string
		Timestamp  int64
		TTLSeconds int32
	}
)

func itemFromProto(itemP *proto.CacheItem) *CacheItem {
	return &CacheItem{
		Value: itemP.Value,
		TTLSeconds: itemP.Ttl,
	}
}