package strato

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestClient(t *testing.T) {
	is := assert.New(t)

	srv, err := NewServer(goodServerCfg)
	is.NoError(err)

	go func() {
		is.NoError(srv.Start())

		defer srv.ShutDown()
	}()

	cl, err := NewClient(goodClientCfg)

	t.Run("Instantiation", func(t *testing.T) {
		is.NoError(err)
		is.NotNil(cl)
	})

	t.Run("KV", func(t *testing.T) {
		goodLoc := &Location{
			Key: "exists",
		}

		val := &Value{
			Content: []byte("some test content"),
		}

		err := cl.Put(goodLoc, val)
		is.NoError(err)

		fetched, err := cl.Get(goodLoc)
		is.NoError(err)
		is.NotNil(fetched)

		badLoc := &Location{
			Key: "does-not-exist",
		}

		fetched, err = cl.Get(badLoc)
		is.Error(err)
		stat, ok := status.FromError(err)
		is.True(ok)
		is.Equal(stat.Code(), codes.NotFound)
		is.Equal(stat.Message(), NotFound(badLoc).Error())
		is.Nil(fetched)
	})
}
