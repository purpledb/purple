package strato

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	is := assert.New(t)

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
		loc := &Location{
			Key: "some-key",
		}

		err := NotFound(loc)
		is.Errorf(err, "not found: no value found for %s", loc.String())
		is.True(IsNotFound(err))
		is.Equal(err.Error(), "no value found for Location<key: some-key>")
	})
}
