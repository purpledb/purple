package grpc

import (
	"context"
	"testing"

	"github.com/lucperkins/strato"

	"github.com/lucperkins/strato/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var goodServerCfg = &strato.ServerConfig{
	Port:    2222,
	Backend: "memory",
}

func TestGrpcServer(t *testing.T) {
	is := assert.New(t)

	ctx := context.Background()

	srv, err := NewGrpcServer(goodServerCfg)
	is.NoError(err)
	is.NoError(srv.backend.Flush())

	go func() {
		is.NoError(srv.Start())
	}()

	t.Run("Instantiation", func(_ *testing.T) {
		is.NoError(err)
		is.NotNil(srv)
	})

	t.Run("Cache", func(_ *testing.T) {
		setReq := &proto.CacheSetRequest{
			Key: "key",
			Item: &proto.CacheItem{
				Value: "value",
				Ttl:   2,
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

	t.Run("Counter", func(_ *testing.T) {
		getReq := &proto.GetCounterRequest{
			Key: "player1",
		}

		res, err := srv.GetCounter(ctx, getReq)
		is.NoError(err)
		is.NotNil(res)
		is.Zero(res.Value)

		incrReq := &proto.IncrementCounterRequest{
			Key:    "player1",
			Amount: 100,
		}

		empty, err := srv.IncrementCounter(ctx, incrReq)
		is.NoError(err)
		is.NotNil(empty)

		res, err = srv.GetCounter(ctx, getReq)
		is.NoError(err)
		is.Equal(res.Value, int64(100))
	})

	t.Run("KV", func(_ *testing.T) {
		locationReq := &proto.Location{
			Key: "key",
		}

		val, err := srv.KVGet(ctx, locationReq)

		stat, ok := status.FromError(err)
		is.True(ok)
		is.Equal(stat.Code(), codes.NotFound)
		is.Nil(val)

		putReq := &proto.PutRequest{
			Location: &proto.Location{
				Key: "key",
			},
			Value: &proto.Value{
				Content: []byte("some content"),
			},
		}

		empty, err := srv.KVPut(ctx, putReq)
		is.NoError(err)
		is.NotNil(empty)

		val, err = srv.KVGet(ctx, locationReq)
		is.NoError(err)
		is.NotNil(val)
		is.Equal(val.Value.Content, []byte("some content"))

		empty, err = srv.KVDelete(ctx, locationReq)
		is.NoError(err)
		is.NotNil(empty)

		val, err = srv.KVGet(ctx, locationReq)
		stat, ok = status.FromError(err)
		is.True(ok)
		is.Equal(stat.Code(), codes.NotFound)
	})

	t.Run("Set", func(_ *testing.T) {
		getReq := &proto.GetSetRequest{
			Set: "set1",
		}

		set, err := srv.SetGet(ctx, getReq)
		is.Nil(err)
		is.Empty(set.Items)

		modifyReq := &proto.ModifySetRequest{
			Set:  "set1",
			Item: "item1",
		}

		empty, err := srv.SetAdd(ctx, modifyReq)
		is.NoError(err)
		is.NotNil(empty)

		set, err = srv.SetGet(ctx, getReq)
		is.NoError(err)
		is.Equal(set.Items, []string{"item1"})

		empty, err = srv.SetRemove(ctx, modifyReq)
		is.NoError(err)
		is.NotNil(empty)

		set, err = srv.SetGet(ctx, getReq)
		is.NoError(err)
		is.Equal(set.Items, []string{})
	})

	t.Run("Shutdown", func(_ *testing.T) {
		is.NoError(srv.ShutDown())
	})
}
