package purple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcClient(t *testing.T) {
	is := assert.New(t)

	cl, err := NewGrpcClient(goodClientCfg)

	t.Run("Instantiation", func(t *testing.T) {
		is.NoError(err)
		is.NotNil(cl)

		noAddressCfg := &ClientConfig{
			Address: "",
		}

		noClient, err := NewGrpcClient(noAddressCfg)
		is.Error(err, ErrNoAddress)
		is.Nil(noClient)

		badAddressCfg := &ClientConfig{
			Address: "1:2:3",
		}
		badCl, err := NewGrpcClient(badAddressCfg)
		is.NoError(err)
		is.NotNil(badCl)

		err = badCl.KVDelete("does-not-exist")
		stat, ok := status.FromError(err)
		is.True(ok)
		is.Equal(stat.Code(), codes.Unavailable)
	})
}
