package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	goodCfg = &Config{
		Port: 2222,
	}

	badCfg = &Config{}
)

func TestConfigInstantiation(t *testing.T) {
	is := assert.New(t)

	err := badCfg.validate()
	is.True(IsConfigError(err))

	is.NoError(goodCfg.validate())
}
