package backend

import (
	"fmt"
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend/disk"
	"github.com/lucperkins/strato/internal/backend/memory"
	"github.com/lucperkins/strato/internal/services/kv"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestServices(t *testing.T) {
	svcs := getServices(t)

	for _, svc := range svcs {
		testSvc(svc, t)
	}
}

func getServices(t *testing.T) []Service {
	is := assert.New(t)

	mem := memory.NewMemoryBackend()

	ds, err := disk.NewDiskBackend()
	is.NoError(err)
	is.NotNil(ds)

	return []Service{mem, ds}
}

func testSvc(svc Service, t *testing.T) {
	is := assert.New(t)

	t.Run(fmt.Sprintf("%s/%s", strings.Title(svc.Name()), "Cache"), func(t *testing.T) {
		is.NoError(svc.Flush())

		key, value := "some-key", "some-value"

		is.NoError(svc.CacheSet(key, value, int32(5)))

		val, err := svc.CacheGet(key)
		is.NoError(err)
		is.NotNil(val)
		is.Equal(val, value)

		is.NoError(svc.CacheSet(key, value, int32(1)))
		time.Sleep(2 * time.Second)
		val, err = svc.CacheGet(key)
		is.True(strato.IsNotFound(err))
		is.Empty(val)

		val, err = svc.CacheGet("does-not-exist")
		is.True(strato.IsNotFound(err))
		is.Empty(val)

		err = svc.CacheSet("", "something", 5)
		is.Equal(err, strato.ErrNoKey)

		err = svc.CacheSet("some-key", "", 5)
		is.Equal(err, strato.ErrNoValue)

		is.NoError(svc.Flush())
	})

	t.Run(fmt.Sprintf("%s/%s", strings.Title(svc.Name()), "Counter"), func(t *testing.T) {
		is.NoError(svc.Flush())

		key, incr := "my-counter", int64(10)

		is.Zero(svc.CounterGet(key))

		is.NoError(svc.CounterIncrement(key, incr))

		val, err := svc.CounterGet(key)
		is.NoError(err)
		is.Equal(val, incr)

		is.NoError(svc.CounterIncrement(key, int64(-50)))
		val, err = svc.CounterGet(key)
		is.NoError(err)
		is.Equal(val, int64(-40))

		val, err = svc.CounterGet("does-not-yet-exist")
		is.NoError(err)
		is.Zero(val)

		is.NoError(svc.Flush())
	})

	t.Run(fmt.Sprintf("%s/%s", strings.Title(svc.Name()), "Flag"), func(t *testing.T) {
		is.NoError(svc.Flush())

		key := "some-flag"

		val, err := svc.FlagGet(key)
		is.NoError(err)
		is.False(val)

		is.NoError(svc.FlagSet(key))

		val, err = svc.FlagGet(key)
		is.NoError(err)
		is.True(val)

		is.NoError(svc.Flush())
	})

	t.Run(fmt.Sprintf("%s/%s", strings.Title(svc.Name()), "KV"), func(t *testing.T) {
		is.NoError(svc.Flush())

		key := "some-key"

		val := &kv.Value{
			Content: []byte("here is a value"),
		}

		is.NoError(svc.KVPut(key, val))

		fetched, err := svc.KVGet("does-not-exist")
		is.True(strato.IsNotFound(err))
		is.Nil(fetched)

		fetched, err = svc.KVGet(key)
		is.NoError(err)
		is.NotNil(fetched)
		is.Equal(fetched, val)

		is.NoError(svc.KVDelete(key))
		fetched, err = svc.KVGet(key)
		is.True(strato.IsNotFound(err))
		is.Nil(fetched)

		is.NoError(svc.Flush())
	})

	t.Run(fmt.Sprintf("%s/%s", strings.Title(svc.Name()), "Set"), func(t *testing.T) {
		is.NoError(svc.Flush())

		set, item1, item2 := "example-set", "example-item-1", "example-item-2"

		is.Empty(svc.SetGet(set))

		s, err := svc.SetAdd(set, item1)
		is.NoError(err)
		is.Len(s, 1)
		is.Equal(s[0], item1)

		is.NotEmpty(svc.SetGet(set))

		s, err = svc.SetAdd(set, item2)
		is.NoError(err)
		is.Len(s, 2)

		s, err = svc.SetGet(set)
		is.Len(s, 2)

		s, err = svc.SetRemove(set, item1)
		is.NoError(err)
		is.Len(s, 1)
		is.Equal(s[0], item2)

		s, err = svc.SetGet(set)
		is.NoError(err)
		is.Len(s, 1)
		is.Equal(s[0], item2)

		s, err = svc.SetRemove(set, item2)
		is.NoError(err)
		is.Empty(s)

		is.NoError(svc.Flush())
	})
}
