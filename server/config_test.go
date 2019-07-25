package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	goodCfg = &Config{
		Port: 2222,
	}

	badCfgNoPort = &Config{}

	badCfgRange = &Config{
		Port: 10,
	}
)

func TestConfigInstantiation(t *testing.T) {
	is := assert.New(t)

	err := badCfgNoPort.validate()
	is.True(IsConfigError(err))
	is.Equal(err, ErrNoPort)

	err = badCfgRange.validate()
	is.True(IsConfigError(err))
	is.Equal(err, ErrPortOutOfRange)

	is.NoError(goodCfg.validate())

	srv, err := New(badCfgNoPort)
	is.True(IsConfigError(err))
	is.Error(err, ErrNoPort)
	is.Nil(srv)

	srv, err = New(goodCfg)
	is.NoError(err)
	is.NotNil(srv)
}
