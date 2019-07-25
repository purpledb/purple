package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorLogic(t *testing.T) {
	is := assert.New(t)

	err := ConfigError{"something went wrong"}
	is.Equal(err.Error(), "server config error: something went wrong")
}
