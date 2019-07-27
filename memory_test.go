package strato

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryImpl(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	mem := New()

	t.Run("Instantiation", func(t *testing.T) {
		is.NotNil(mem)
		is.Empty(mem.values)
	})

	t.Run("Cache", func(t *testing.T) {
		key, value := "some-key", "some-value"

		is.NoError(mem.CacheSet(key, value, 5))
		val, err := mem.CacheGet(key)
		is.NoError(err)
		is.Equal(val, value)

		is.NoError(mem.CacheSet(key, value, 1))
		time.Sleep(2 * time.Second)
		val, err = mem.CacheGet(key)
		is.True(IsExpired(err))
		is.Empty(val)

		val, err = mem.CacheGet("does-not-exist")
		is.True(IsNoItemFound(err))
		is.Empty(val)

		err = mem.CacheSet("", "something", 5)
		is.True(IsNoCacheKey(err))
		err = mem.CacheSet("some-key", "", 5)
		is.True(IsNoCacheValue(err))
	})

	t.Run("Counter", func(t *testing.T) {
		key := "my-counter"

		is.Zero(mem.GetCounter(key))

		mem.IncrementCounter(key, int32(10))
		is.Equal(mem.GetCounter(key), int32(10))
		mem.IncrementCounter(key, int32(-50))
		is.Equal(mem.GetCounter(key), int32(-40))
		is.Zero(mem.GetCounter("does-not-yet-exist"), 0)
	})

	t.Run("KV", func(t *testing.T) {
		loc := &Location{
			Key: "some-key",
		}

		val := &Value{
			Content: []byte("here is a value"),
		}

		mem.KVPut(loc, val)

		fetched, err := mem.KVGet(&Location{Key: "does-not-exist"})
		is.True(IsNotFound(err))
		is.Nil(fetched)

		fetched, err = mem.KVGet(loc)
		is.NoError(err)
		is.NotNil(fetched)
		is.Equal(fetched, val)

		mem.KVDelete(loc)
		fetched, err = mem.KVGet(loc)
		is.True(IsNotFound(err))
		is.Nil(fetched)
	})

	t.Run("Search", func(t *testing.T) {
		doc := &Document{
			ID:      "doc-1",
			Content: "Here lies searchable content",
		}

		goodQuery, badQuery := "here", "oops"

		res := mem.Query(goodQuery)
		is.Empty(res)

		mem.Index(doc)
		res = mem.Query(goodQuery)
		is.Len(res, 1)

		res = mem.Query(badQuery)
		is.Empty(res)
	})
}
