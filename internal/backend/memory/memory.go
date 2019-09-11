package memory

import (
	"github.com/lucperkins/strato/internal/data"
	"time"

	"github.com/lucperkins/strato/internal/services/cache"
	"github.com/lucperkins/strato/internal/services/counter"
	"github.com/lucperkins/strato/internal/services/kv"
	"github.com/lucperkins/strato/internal/services/set"

	"github.com/lucperkins/strato"
)

type Memory struct {
	cache    map[string]*cache.Item
	counters map[string]int64
	kv       map[string]*kv.Value
	sets     map[string]*data.Set
}

var (
	_ cache.Cache     = (*Memory)(nil)
	_ counter.Counter = (*Memory)(nil)
	_ kv.KV           = (*Memory)(nil)
	_ set.Set         = (*Memory)(nil)
)

func NewMemoryBackend() *Memory {
	cacheMem := make(map[string]*cache.Item)

	counterMem := make(map[string]int64)

	setMem := make(map[string]*data.Set)

	kvMem := make(map[string]*kv.Value)

	return &Memory{
		cache:    cacheMem,
		counters: counterMem,
		kv:       kvMem,
		sets:     setMem,
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
		return "", strato.NotFound(key)
	}

	now := time.Now().Unix()

	expired := (now - val.Timestamp) > int64(val.TTLSeconds)

	if expired {
		delete(m.cache, key)

		return "", strato.NotFound(key)
	}

	return val.Value, nil
}

func (m *Memory) CacheSet(key, value string, ttl int32) error {
	if key == "" {
		return strato.ErrNoKey
	}

	if value == "" {
		return strato.ErrNoValue
	}

	item := &cache.Item{
		Value:      value,
		Timestamp:  time.Now().Unix(),
		TTLSeconds: parseTtl(ttl),
	}

	m.cache[key] = item

	return nil
}

func parseTtl(ttl int32) int32 {
	if ttl == 0 {
		return cache.DefaultTtl
	} else {
		return ttl
	}
}

// Counter
func (m *Memory) CounterIncrement(key string, increment int64) error {
	count, ok := m.counters[key]
	if !ok {
		m.counters[key] = increment
	} else {
		m.counters[key] = count + increment
	}

	return nil
}

func (m *Memory) CounterGet(key string) (int64, error) {
	return m.counters[key], nil
}

func (m *Memory) KVGet(key string) (*kv.Value, error) {
	val, ok := m.kv[key]
	if !ok {
		return nil, strato.NotFound(key)
	}

	return val, nil
}

func (m *Memory) KVPut(key string, value *kv.Value) error {
	m.kv[key] = value
	return nil
}

func (m *Memory) KVDelete(key string) error {
	delete(m.kv, key)
	return nil
}

func (m *Memory) SetGet(set string) ([]string, error) {
	s, ok := m.sets[set]

	if !ok {
		return data.EmptySet(), nil
	}

	return s.Get(), nil
}

func (m *Memory) SetAdd(set, item string) ([]string, error) {
	var result []string

	s, ok := m.sets[set]

	if ok {
		s.Add(item)
		result = s.Get()
	} else {
		newSet := data.NewSet(item)

		m.sets[set] = newSet
		result = newSet.Get()
	}

	return result, nil
}

func (m *Memory) SetRemove(set, item string) ([]string, error) {
	s, ok := m.sets[set]

	if ok {
		s.Remove(item)
		return s.Get(), nil
	} else {
		return data.EmptySet(), nil
	}
}
