package redis

import (
	"github.com/go-redis/redis"
	"github.com/lucperkins/strato"
	"time"
)

const defaultUrl = "localhost:6379"

type Redis struct {
	cache, counters, kv, sets *redis.Client
}

var (
	_ strato.Cache   = (*Redis)(nil)
	_ strato.Counter = (*Redis)(nil)
	_ strato.KV      = (*Redis)(nil)
	_ strato.Set     = (*Redis)(nil)
)

func NewRedisBackend(addr string) (*Redis, error) {
	if addr == "" {
		addr = defaultUrl
	}

	cache, err := newRedisClient(addr, 0)
	if err != nil {
		return nil, err
	}

	counters, err := newRedisClient(addr, 1)
	if err != nil {
		return nil, err
	}

	kv, err := newRedisClient(addr, 2)
	if err != nil {
		return nil, err
	}

	sets, err := newRedisClient(addr, 3)
	if err != nil {
		return nil, err
	}

	return &Redis{
		cache:    cache,
		counters: counters,
		kv:       kv,
		sets:     sets,
	}, nil
}

func newRedisClient(addr string, i int) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: "",
		DB:       i,
	}

	cl := redis.NewClient(opts)

	if err := cl.Ping().Err(); err != nil {
		return nil, err
	}

	return cl, nil
}

// Interface methods
func (r *Redis) Close() error {
	for _, db := range []*redis.Client{
		r.cache,
	} {
		if err := db.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) Flush() error {
	for _, db := range []*redis.Client{
		r.cache,
	} {
		if err := db.FlushAll().Err(); err != nil {
			return err
		}
	}

	return nil
}

// Cache operations

func (r *Redis) CacheGet(key string) (string, error) {
	s, err := r.cache.Get(key).Result()

	if err != nil {
		if err == redis.Nil {
			return "", strato.NotFound(key)
		} else {
			return "", err
		}
	}

	return s, nil
}

func (r *Redis) CacheSet(key, value string, ttl int32) error {
	t := time.Duration(ttl) * time.Second

	return r.cache.Set(key, value, t).Err()
}

// Counter operations

func (r *Redis) CounterGet(key string) (int64, error) {
	i, err := r.counters.Get(key).Int64()

	if err == redis.Nil {
		return 0, nil
	} else {
		return i, err
	}
}

func (r *Redis) CounterIncrement(key string, increment int64) error {
	return r.counters.IncrBy(key, increment).Err()
}

// KV operations

func (r *Redis) KVGet(key string) (*strato.Value, error) {
	s, err := r.kv.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, strato.NotFound(key)
		} else {
			return nil, err
		}
	}

	return &strato.Value{
		Content: []byte(s),
	}, nil
}

func (r *Redis) KVPut(key string, value *strato.Value) error {
	return r.kv.Set(key, value.Content, 0).Err()
}

func (r *Redis) KVDelete(key string) error {
	return r.kv.Del(key).Err()
}

// Set operations

func (r *Redis) GetSet(set string) ([]string, error) {
	return r.sets.SMembers(set).Result()
}

func (r *Redis) AddToSet(set, item string) error {
	return r.sets.SAdd(set, item).Err()
}

func (r *Redis) RemoveFromSet(set, item string) error {
	return r.sets.SRem(set, item).Err()
}
