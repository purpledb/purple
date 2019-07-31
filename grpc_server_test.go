package strato

import (
	"context"
	"github.com/lucperkins/strato/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrpcServer(t *testing.T) {
	is := assert.New(t)

	ctx := context.Background()

	srv, err := NewGrpcServer(goodServerCfg)

	t.Run("Instantiation", func(_ *testing.T) {
		is.NoError(err)
		is.NotNil(srv)
	})

	t.Run("Cache", func(_ *testing.T) {
		setReq := &proto.CacheSetRequest{
			Key: "key",
			Item: &proto.CacheItem{
				Value: "value",
				Ttl: 2,
			},
		}

		empty, err := srv.CacheSet(ctx, setReq)
		is.NoError(err)
		is.NotNil(empty)

		getReq := &proto.CacheGetRequest{
			Key: "key",
		}

		res, err := srv.CacheGet(ctx, getReq)
		is.NoError(err)
		is.NotNil(res)
		is.Equal(res.Value, "value")
	})

	is.NoError(srv.mem.Close())
}
