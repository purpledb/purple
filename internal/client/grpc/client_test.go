package grpc

import (
	"github.com/lucperkins/strato/internal/config"
	"testing"

	"github.com/lucperkins/strato"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var goodClientCfg = &config.ClientConfig{
	Address: "localhost:2222",
}

func TestGrpcClient(t *testing.T) {
	is := assert.New(t)

	cl, err := NewClient(goodClientCfg)

	t.Run("Instantiation", func(t *testing.T) {
		is.NoError(err)
		is.NotNil(cl)

		noAddressCfg := &config.ClientConfig{
			Address: "",
		}

		noClient, err := NewClient(noAddressCfg)
		is.Error(err, strato.ErrNoAddress)
		is.Nil(noClient)

		badAddressCfg := &config.ClientConfig{
			Address: "1:2:3",
		}
		badCl, err := NewClient(badAddressCfg)
		is.NoError(err)
		is.NotNil(badCl)

		err = badCl.KVDelete("does-not-exist")
		stat, ok := status.FromError(err)
		is.True(ok)
		is.Equal(stat.Code(), codes.Unavailable)
	})
}
