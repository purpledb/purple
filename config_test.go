package strato

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	goodClientCfg = &ClientConfig{
		Address: "localhost:2222",
	}

	goodServerCfg = &ServerConfig{
		Port: 2222,
	}
)

func TestConfigInstantiation(t *testing.T) {
	is := assert.New(t)

	t.Run("Client", func(t *testing.T) {
		emptyCfg := &ClientConfig{}

		err := emptyCfg.validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrNoAddress)

		err = goodClientCfg.validate()
		is.NoError(err)
	})

	t.Run("Server", func(t *testing.T) {
		emptyCfg := &ServerConfig{}

		lowPortCfg := &ServerConfig{
			Port: 10,
		}

		err := emptyCfg.validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrNoPort)

		err = lowPortCfg.validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrPortOutOfRange)

		is.NoError(goodServerCfg.validate())

		srv, err := NewServer(emptyCfg)
		is.True(IsConfigError(err))
		is.Equal(err, ErrNoPort)
		is.Nil(srv)

		srv, err = NewServer(lowPortCfg)
		is.True(IsConfigError(err))
		is.Equal(err, ErrPortOutOfRange)

		srv, err = NewServer(goodServerCfg)
		is.NoError(err)
		is.NotNil(srv)
	})
}
