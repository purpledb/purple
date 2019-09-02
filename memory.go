package strato

import (
	bolt "github.com/etcd-io/bbolt"
	"os"
	"time"
)

const dbFile = "strato-kv.db"

type Memory struct {
	cache    map[string]*CacheItem
	counters map[string]int64
	kv       *bolt.DB
	sets     map[string][]string
}

var (
	_ Cache   = (*Memory)(nil)
	_ Counter = (*Memory)(nil)
	_ KV      = (*Memory)(nil)
	_ Set     = (*Memory)(nil)
)

func NewMemoryBackend() *Memory {
	cache := make(map[string]*CacheItem)

	counters := make(map[string]int64)

	sets := make(map[string][]string)

	if _, err := os.Create(dbFile); err != nil {
		panic(err)
	}

	kv, err := bolt.Open(dbFile, 0666, nil)
	if err != nil {
		panic(err)
	}

	return &Memory{
		cache:    cache,
		counters: counters,
		kv:       kv,
		sets:     sets,
	}
}

// Backend methods
func (m *Memory) Close() error {
	return m.kv.Close()
}

// Cache
func (m *Memory) CacheGet(key string) (string, error) {
	val, ok := m.cache[key]

	if !ok {
		return "", ErrNoCacheItem
	}

	now := time.Now().Unix()

	expired := (now - val.Timestamp) > int64(val.TTLSeconds)

	if expired {
		delete(m.cache, key)

		return "", ErrExpired
	}

	return val.Value, nil
}

func (m *Memory) CacheSet(key, value string, ttl int32) error {
	if key == "" {
		return ErrNoCacheKey
	}

	if value == "" {
		return ErrNoCacheValue
	}

	item := &CacheItem{
		Value:      value,
		Timestamp:  time.Now().Unix(),
		TTLSeconds: getTtl(ttl),
	}

	m.cache[key] = item

	return nil
}

func getTtl(ttl int32) int32 {
	if ttl == 0 {
		return defaultTtl
	} else {
		return ttl
	}
}

// Counter
func (m *Memory) CounterIncrement(key string, increment int64) error {
	counter, ok := m.counters[key]
	if !ok {
		m.counters[key] = increment
	} else {
		m.counters[key] = counter + increment
	}

	return nil
}

func (m *Memory) CounterGet(key string) (int64, error) {
	return m.counters[key], nil
}

func (m *Memory) KVGet(location *Location) (*Value, error) {
	if err := location.validate(); err != nil {
		return nil, err
	}

	var val []byte

	if err := m.kv.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(location.Bucket))

		if buck == nil {
			return NotFound(location)
		}

		val = buck.Get([]byte(location.Key))

		if val == nil {
			return NotFound(location)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Value{
		Content: val,
	}, nil
}

func (m *Memory) KVPut(location *Location, value *Value) error {
	if err := location.validate(); err != nil {
		return err
	}

	return m.kv.Update(func(tx *bolt.Tx) error {
		buck, err := tx.CreateBucketIfNotExists([]byte(location.Bucket))
		if err != nil {
			return err
		}

		return buck.Put([]byte(location.Key), []byte(value.Content))
	})
}

func (m *Memory) KVDelete(location *Location) error {
	if err := location.validate(); err != nil {
		return err
	}

	return m.kv.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(location.Bucket))
		if buck == nil {
			return nil
		}

		return buck.Delete([]byte(location.Key))
	})
}

func (m *Memory) GetSet(set string) ([]string, error) {
	s, ok := m.sets[set]

	if !ok {
		return []string{}, nil
	}

	return s, nil
}

func (m *Memory) AddToSet(set, item string) error {
	if _, ok := m.sets[set]; ok {
		m.sets[set] = append(m.sets[set], item)
	} else {
		m.sets[set] = []string{item}
	}

	return nil
}

func (m *Memory) RemoveFromSet(set, item string) error {
	_, ok := m.sets[set]
	if ok {
		for idx, it := range m.sets[set] {
			if it == item {
				m.sets[set] = append(m.sets[set][:idx], m.sets[set][idx+1:]...)
			}
		}

		return nil
	} else {
		return ErrNoSet
	}
}
