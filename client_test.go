package strato

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		noAddressCfg := &ClientConfig{
			Address: "",
		}

		noClient, err := NewClient(noAddressCfg)
		is.Error(err, ErrNoAddress)
		is.Nil(noClient)

		badAddressCfg := &ClientConfig{
			Address: "1:2:3",
		}
		badCl, err := NewClient(badAddressCfg)
		is.NoError(err)
		is.NotNil(badCl)

		err = badCl.Delete(&Location{Key: "does-not-exist"})
		stat, ok := status.FromError(err)
		is.True(ok)
		is.Equal(stat.Code(), codes.Unavailable)
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
		stat := status.Convert(err)
		is.Equal(stat.Code(), codes.NotFound)
		is.Equal(stat.Message(), NotFound(badLoc).Error())
		is.Nil(fetched)

		t.Run("Nils", func(t *testing.T) {
			err = cl.Put(nil, nil)
			is.Equal(err, ErrNoLocation)
			err = cl.Put(&Location{Key: "test"}, nil)
			is.Equal(err, ErrNoValue)
			err = cl.Put(nil, &Value{Content: []byte("some bytes")})
		})
	})

	t.Run("Search", func(t *testing.T) {
		doc := &Document{
			ID:      "doc-100",
			Content: "This is the 100th DOC",
		}

		goodQ, badQ := "doc", "does not exist"

		res, err := cl.Query(goodQ)

		is.NoError(err)
		is.Empty(res)

		is.NoError(cl.Index(doc))

		res, err = cl.Query(badQ)
		is.NoError(err)
		is.Empty(res)

		res, err = cl.Query(goodQ)
		is.NoError(err)
		is.Len(res, 1)
	})
}
