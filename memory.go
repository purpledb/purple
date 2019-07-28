package strato

import (
	"strings"
	"time"
)

type Memory struct {
	cache    map[string]*CacheItem
	counters map[string]int32
	values   map[Location]*Value
	docs     []*Document
	sets     map[string][]string
}

var (
	_ Cache   = (*Memory)(nil)
	_ Counter = (*Memory)(nil)
	_ KV      = (*Memory)(nil)
	_ Search  = (*Memory)(nil)
	_ Set     = (*Memory)(nil)
)

func NewMemory() *Memory {
	cache := make(map[string]*CacheItem)

	counters := make(map[string]int32)

	values := make(map[Location]*Value)

	docs := make([]*Document, 0)

	sets := make(map[string][]string)

	return &Memory{
		cache:    cache,
		counters: counters,
		values:   values,
		docs:     docs,
		sets:     sets,
	}
}

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

func (m *Memory) IncrementCounter(key string, increment int32) {
	counter, ok := m.counters[key]
	if !ok {
		m.counters[key] = increment
	} else {
		m.counters[key] = counter + increment
	}
}

func (m *Memory) GetCounter(key string) int32 {
	return m.counters[key]
}

func (m *Memory) KVGet(location *Location) (*Value, error) {
	val, ok := m.values[*location]

	if !ok {
		return nil, NotFound(location)
	}

	return val, nil
}

func (m *Memory) KVPut(location *Location, value *Value) {
	m.values[*location] = value
}

func (m *Memory) KVDelete(location *Location) {
	if _, ok := m.values[*location]; ok {
		delete(m.values, *location)
	}
}

func (m *Memory) Index(doc *Document) {
	doc = doc.prepare()

	m.docs = append(m.docs, doc)
}

func (m *Memory) Query(q string) []*Document {
	strings.ToLower(q)

	docs := make([]*Document, 0)

	if len(m.docs) == 0 {
		return []*Document{}
	}

	for _, doc := range m.docs {
		if strings.Contains(doc.Content, q) {
			docs = append(docs, doc)
		}
	}

	return docs
}

func (m *Memory) GetSet(set string) []string {
	s, ok := m.sets[set]

	if !ok {
		return []string{}
	}

	return s
}

func (m *Memory) AddToSet(set, item string) {
	if _, ok := m.sets[set]; !ok {
		m.sets[set] = []string{item}
		return
	}

	m.sets[set] = append(m.sets[set], item)
}

func (m *Memory) RemoveFromSet(set, item string) {
	if _, ok := m.sets[set]; !ok {
		return
	}

	for idx, it := range m.sets[set] {
		if it == item {
			m.sets[set] = append(m.sets[set][:idx], m.sets[set][idx+1:]...)
		}
	}
}
