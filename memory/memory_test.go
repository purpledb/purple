package memory

import (
	"github.com/stretchr/testify/assert"
	"strato/kv"
	"testing"
)

func TestMemoryImpl(t *testing.T) {
	is := assert.New(t)

	mem := New()
	is.NotNil(mem)
	is.Empty(mem.values)

	loc := &kv.Location{
		Key: "some-key",
	}

	val := &kv.Value{
		Content: []byte("here is a value"),
	}

	is.NoError(mem.Put(loc, val))

	fetched, err := mem.Get(&kv.Location{Key: "does-not-exist"})
	is.True(kv.IsNotFound(err))
	is.Nil(fetched)

	fetched, err = mem.Get(loc)
	is.NoError(err)
	is.NotNil(fetched)

	is.NoError(mem.Delete(loc))
	fetched, err = mem.Get(loc)
	is.True(kv.IsNotFound(err))
	is.Nil(fetched)
}
