package strato

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryImpl(t *testing.T) {
	is := assert.New(t)

	mem := New()

	t.Run("Instantiation", func(t *testing.T) {
		is.NotNil(mem)
		is.Empty(mem.values)
	})

	t.Run("KV", func(t *testing.T) {
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
	})

	t.Run("Search", func(t *testing.T) {
		doc := &Document{
			ID: "doc-1",
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
