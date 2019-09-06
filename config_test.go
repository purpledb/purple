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

		err := emptyCfg.Validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrNoAddress)

		err = goodClientCfg.Validate()
		is.NoError(err)
	})

	t.Run("GrpcServer", func(t *testing.T) {
		emptyCfg := &ServerConfig{}

		lowPortCfg := &ServerConfig{
			Port: 10,
		}

		err := emptyCfg.Validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrNoPort)

		err = lowPortCfg.Validate()
		is.True(IsConfigError(err))
		is.Equal(err, ErrPortOutOfRange)

		is.NoError(goodServerCfg.Validate())
	})
}
