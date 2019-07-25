package strato

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strato/proto"
	"testing"
)

const (
	goodKey = "exists"
	badKey  = "does-not-exist"
)

var (
	ctx = context.Background()

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

	srv, err := NewServer(goodServerCfg)
	is.NoError(err)
	is.NotNil(srv)

	t.Run("KV", func(t *testing.T) {
		empty, err := srv.Put(ctx, goodReq)
		is.NoError(err)
		is.NotNil(empty)

		fetched, err := srv.Get(ctx, goodLoc)
		is.NoError(err)
		is.NotNil(fetched)
		is.Equal(fetched.Value.Content, goodVal.Content)

		empty, err = srv.Delete(ctx, goodLoc)
		is.NoError(err)
		is.NotNil(empty)

		fetched, err = srv.Get(ctx, badLoc)
		is.True(IsNotFound(err))
		is.Nil(fetched)
	})

	t.Run("Start/Shutdown", func(t *testing.T) {
		go func() {
			is.NoError(srv.Start())
			srv.ShutDown()
		}()
	})
}
