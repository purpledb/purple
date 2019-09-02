package strato

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	goodClientCfg = &ClientConfig{
		Address: "localhost:2222",
	}

	goodServerCfg = &ServerConfig{
		Port:    2222,
		Backend: "disk",
	}
)

func TestConfigInstantiation(t *testing.T) {
	is := assert.New(t)

	t.Run("GrpcClient", func(t *testing.T) {
		emptyCfg := &ClientConfig{}

		err := emptyCfg.validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrNoAddress)

		err = goodClientCfg.validate()
		is.NoError(err)
	})

	t.Run("GrpcServer", func(t *testing.T) {
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
	})
}
