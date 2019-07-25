package strato

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	goodCfg = &ServerConfig{
		Port: 2222,
	}

	badCfgNoPort = &ServerConfig{}

	badCfgRange = &ServerConfig{
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

	srv, err := NewServer(badCfgNoPort)
	is.True(IsConfigError(err))
	is.Error(err, ErrNoPort)
	is.Nil(srv)

	srv, err = NewServer(goodCfg)
	is.NoError(err)
	is.NotNil(srv)
}
