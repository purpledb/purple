package memory

import (
	"time"

	"github.com/lucperkins/strato"
)

type Memory struct {
	cache    map[string]*strato.CacheItem
	counters map[string]int64
	kv       map[string]*strato.Value
	sets     map[string][]string
}

var (
	_ strato.Cache   = (*Memory)(nil)
	_ strato.Counter = (*Memory)(nil)
	_ strato.KV      = (*Memory)(nil)
	_ strato.Set     = (*Memory)(nil)
)

func NewMemoryBackend() *Memory {
	cache := make(map[string]*strato.CacheItem)

	counters := make(map[string]int64)

	sets := make(map[string][]string)

	kv := make(map[string]*strato.Value)

	return &Memory{
		cache:    cache,
		counters: counters,
		kv:       kv,
		sets:     sets,
	}
}

// Interface methods
func (m *Memory) Close() error {
	return nil
}

func (m *Memory) Flush() error {
	return nil
}

// Cache
func (m *Memory) CacheGet(key string) (string, error) {
	val, ok := m.cache[key]

	if !ok {
		return "", strato.ErrNoCacheItem
	}

	now := time.Now().Unix()

	expired := (now - val.Timestamp) > int64(val.TTLSeconds)

	if expired {
		delete(m.cache, key)

		return "", strato.ErrExpired
	}

	return val.Value, nil
}

func (m *Memory) CacheSet(key, value string, ttl int32) error {
	if key == "" {
		return strato.ErrNoCacheKey
	}

	if value == "" {
		return strato.ErrNoCacheValue
	}

	item := &strato.CacheItem{
		Value:      value,
		Timestamp:  time.Now().Unix(),
		TTLSeconds: parseTtl(ttl),
	}

	m.cache[key] = item

	return nil
}

func parseTtl(ttl int32) int32 {
	if ttl == 0 {
		return strato.DefaultTtl
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

func (m *Memory) KVGet(key string) (*strato.Value, error) {
	val, ok := m.kv[key]
	if !ok {
		return nil, strato.NotFound(key)
	}

	return val, nil
}

func (m *Memory) KVPut(key string, value *strato.Value) error {
	m.kv[key] = value
	return nil
}

func (m *Memory) KVDelete(key string) error {
	delete(m.kv, key)
	return nil
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
		return strato.ErrNoSet
	}
}
