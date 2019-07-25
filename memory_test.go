package strato

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryImpl(t *testing.T) {
	is := assert.New(t)

	mem := New()
	is.NotNil(mem)
	is.Empty(mem.values)

	loc := &Location{
		Key: "some-key",
	}

	val := &Value{
		Content: []byte("here is a value"),
	}

	mem.Put(loc, val)

	fetched, err := mem.Get(&Location{Key: "does-not-exist"})
	is.True(IsNotFound(err))
	is.Nil(fetched)

	fetched, err = mem.Get(loc)
	is.NoError(err)
	is.NotNil(fetched)
	is.Equal(fetched, val)

	mem.Delete(loc)
	fetched, err = mem.Get(loc)
	is.True(IsNotFound(err))
	is.Nil(fetched)
}
