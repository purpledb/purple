package memory

import (
	"testing"
	"time"

	"github.com/lucperkins/strato/internal/services/kv"

	"github.com/lucperkins/strato"

	"github.com/stretchr/testify/assert"
)

func TestMemoryImpl(t *testing.T) {
	is := assert.New(t)

	mem := NewMemoryBackend()

	t.Run("Instantiation", func(t *testing.T) {
		is.NotNil(mem)
		is.NotNil(mem.kv)
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
		is.True(strato.IsNotFound(err))
		is.Empty(val)

		val, err = mem.CacheGet("does-not-exist")
		is.True(strato.IsNotFound(err))
		is.Empty(val)

		err = mem.CacheSet("", "something", 5)
		is.Equal(err, strato.ErrNoKey)

		err = mem.CacheSet("some-key", "", 5)
		is.Equal(err, strato.ErrNoValue)
	})

	t.Run("Counter", func(t *testing.T) {
		key, incr := "my-counter", int64(10)

		is.Zero(mem.CounterGet(key))

		is.NoError(mem.CounterIncrement(key, incr))

		val, err := mem.CounterGet(key)
		is.NoError(err)
		is.Equal(val, incr)

		is.NoError(mem.CounterIncrement(key, int64(-50)))
		val, err = mem.CounterGet(key)
		is.NoError(err)
		is.Equal(val, int64(-40))

		val, err = mem.CounterGet("does-not-yet-exist")
		is.NoError(err)
		is.Zero(val)
	})

	t.Run("KV", func(t *testing.T) {
		key := "some-key"

		val := &kv.Value{
			Content: []byte("here is a value"),
		}

		is.NoError(mem.KVPut(key, val))

		fetched, err := mem.KVGet("does-not-exist")
		is.True(strato.IsNotFound(err))
		is.Nil(fetched)

		fetched, err = mem.KVGet(key)
		is.NoError(err)
		is.NotNil(fetched)
		is.Equal(fetched, val)

		is.NoError(mem.KVDelete(key))
		fetched, err = mem.KVGet(key)
		is.True(strato.IsNotFound(err))
		is.Nil(fetched)
	})

	t.Run("Set", func(t *testing.T) {
		set, item1, item2 := "example-set", "example-item-1", "example-item-2"

		is.Empty(mem.SetGet(set))

		s, err := mem.SetAdd(set, item1)
		is.NoError(err)
		is.Len(s, 1)
		is.Equal(s[0], item1)

		is.NotEmpty(mem.SetGet(set))

		s, err = mem.SetAdd(set, item2)
		is.NoError(err)
		is.Len(s, 2)

		s, err = mem.SetGet(set)
		is.Len(s, 2)

		s, err = mem.SetRemove(set, item1)
		is.NoError(err)
		is.Len(s, 1)
		is.Equal(s[0], item2)

		s, err = mem.SetGet(set)
		is.NoError(err)
		is.Len(s, 1)
		is.Equal(s[0], item2)

		s, err = mem.SetRemove(set, item2)
		is.NoError(err)
		is.Empty(s)
	})
}
