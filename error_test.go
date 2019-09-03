package strato

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	is := assert.New(t)

	t.Run("Cache", func(t *testing.T) {
		err := ErrNoCacheItem
		is.Equal(err.Error(), "cache error: no item found")
		is.True(IsNoItemFound(err))
	})

	t.Run("Config", func(t *testing.T) {
		err := ConfigError{"some config error"}
		is.Equal(err.Error(), "config error: some config error")
		is.True(IsConfigError(err))
	})

	t.Run("KV", func(t *testing.T) {
		err := KVError{"some KV error"}
		is.Equal(err.Error(), "KV error: some KV error")
		is.True(IsKVError(err))
	})

	t.Run("NotFound", func(t *testing.T) {
		key := "some-key"

		err := NotFound(key)
		is.Errorf(err, "not found: no value found for %s", key)
		is.True(IsNotFound(err))
		is.Equal(err.Error(), "no value found for some-key")
	})
}
