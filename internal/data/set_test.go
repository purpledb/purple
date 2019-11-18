package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetType(t *testing.T) {
	is := assert.New(t)

	fruits := NewSet()

	fruits.Remove("some-item")

	fruits.Add("apple")
	is.Len(fruits.items, 1)
	is.Equal(fruits.items[0], "apple")

	fruits.Add("apple")
	is.Len(fruits.items, 1)
	is.Equal(fruits.items[0], "apple")

	fruits.Remove("apple")
	is.Empty(fruits.items)

	sets := make(map[string]*Set)
	is.Len(sets, 0)
}
