package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataHelpers(t *testing.T) {
	is := assert.New(t)

	t.Run("Integers", func(t *testing.T) {
		testCases := []int{
			0, 1, 12345, 987654321,
		}

		for _, tc := range testCases {
			i := int64(tc)

			is.Equal(i, BytesToInt64(Int64ToBytes(i)))
		}
	})

	t.Run("Sets", func(t *testing.T) {
		testCases := [][]string{
			{},
			{"apple", "orange", "banana"},
			{"just-one"},
			{"longer string", "1234567"},
		}

		for _, tc := range testCases {
			s := NewSet(tc...)

			bs, err := s.ToBytes()
			is.NoError(err)
			is.NotNil(bs)

			res, err := BytesToSet(bs)
			is.NoError(err)
			is.Equal(tc, res.items)
		}
	})
}
