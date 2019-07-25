package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strato/kv"
	"strato/proto"
	"testing"
)

const (
	goodKey = "exists"
	badKey  = "does-not-exist"
)

var (
	goodLoc = &proto.Location{
		Key: goodKey,
	}

	badLoc = &proto.Location{
		Key: badKey,
	}

	goodContent = []byte("here is some test value content")

	goodVal = &proto.Value{
		Content: goodContent,
	}

	goodReq = &proto.PutRequest{
		Location: goodLoc,
		Value:    goodVal,
	}
)

func TestServer(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()

	cfg := &Config{
		Port: 8080,
	}

	srv, err := New(cfg)
	is.NoError(err)
	is.NotNil(srv)

	t.Run("Start", func(t *testing.T) {
		go func() {
			is.NoError(srv.Start())
		}()

		srv.ShutDown()
	})

	t.Run("KV", func(t *testing.T) {
		empty, err := srv.Put(ctx, goodReq)
		is.NoError(err)
		is.NotNil(empty)
		is.NotEmpty(srv.mem.All())

		fetched, err := srv.Get(ctx, goodLoc)
		is.NoError(err)
		is.NotNil(fetched)
		is.Equal(fetched.Value.Content, goodVal.Content)

		empty, err = srv.Delete(ctx, goodLoc)
		is.NoError(err)
		is.NotNil(empty)

		fetched, err = srv.Get(ctx, badLoc)
		is.True(kv.IsNotFound(err))
		is.Nil(fetched)
	})
}
