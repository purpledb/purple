package config

import (
	"github.com/lucperkins/strato"
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
		is.True(strato.IsConfigError(err))
		is.Equal(err, strato.ErrNoAddress)

		err = goodClientCfg.Validate()
		is.NoError(err)
	})

	t.Run("GrpcServer", func(t *testing.T) {
		emptyCfg := &ServerConfig{}

		lowPortCfg := &ServerConfig{
			Port: 10,
		}

		err := emptyCfg.Validate()
		is.True(strato.IsConfigError(err))
		is.Equal(err, strato.ErrNoPort)

		err = lowPortCfg.Validate()
		is.True(strato.IsConfigError(err))
		is.Equal(err, strato.ErrPortOutOfRange)

		is.NoError(goodServerCfg.Validate())
	})
}
