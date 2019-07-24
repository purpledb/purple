package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigInstantiation(t *testing.T) {
	is := assert.New(t)

	cfg := &Config{}
	err := cfg.validate()
	is.True(IsConfigError(err))

	cfg = &Config{
		Port: 2222,
	}
	is.NoError(cfg.validate())
}
