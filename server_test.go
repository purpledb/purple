package strato

import (
	"context"
	"testing"
	"time"

	"github.com/lucperkins/strato/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	t.Run("Cache", func(t *testing.T) {
		key, value, ttl := "some-key", "some-val", 1
		setReq := &proto.CacheSetRequest{
			Key: key,
			Item: &proto.CacheItem{
				Value: value,
				Ttl:   int32(ttl),
			},
		}

		getReq := &proto.CacheGetRequest{
			Key: key,
		}

		empty, err := srv.CacheSet(ctx, setReq)
		is.NoError(err)
		is.NotNil(empty)

		val, err := srv.CacheGet(ctx, getReq)
		is.NoError(err)
		is.Equal(val.Value, value)

		time.Sleep(2 * time.Second)
		val, err = srv.CacheGet(ctx, getReq)
		is.True(IsExpired(err))
		is.Nil(val)

		badGetReq := &proto.CacheGetRequest{
			Key: "does-not-exist",
		}

		val, err = srv.CacheGet(ctx, badGetReq)
		is.True(IsNoItemFound(err))
		is.Nil(val)

		badSetReq := &proto.CacheSetRequest{
			Key: key,
			Item: &proto.CacheItem{
				Value: "",
				Ttl:   5,
			},
		}

		empty, err = srv.CacheSet(ctx, badSetReq)
		stat := status.Convert(err)
		is.Equal(stat.Code(), codes.Unknown)
		is.Equal(stat.Message(), ErrNoCacheValue.Error())

		is.Nil(empty)
	})

	t.Run("Counter", func(t *testing.T) {
		key, incr := "example-key", int32(25)

		getReq := &proto.GetCounterRequest{
			Key: key,
		}

		res, err := srv.GetCounter(ctx, getReq)
		is.NoError(err)
		is.NotNil(res)
		is.Zero(res.Value)

		empty, err := srv.IncrementCounter(ctx, &proto.IncrementCounterRequest{Key: key, Amount: incr})
		is.NoError(err)
		is.NotNil(empty)

		res, err = srv.GetCounter(ctx, getReq)
		is.NoError(err)
		is.NotNil(res)
		is.Equal(res.Value, incr)
	})

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

		stat := status.Convert(err)
		is.Equal(stat.Code(), codes.NotFound)
		is.Equal(stat.Message(), NotFound(&Location{Key: badKey}).Error())
		is.Nil(fetched)
	})

	t.Run("Search", func(t *testing.T) {
		doc := &proto.Document{
			Id:      "some-id",
			Content: "some content to be searched",
		}

		req := &proto.IndexRequest{
			Document: doc,
		}

		empty, err := srv.Index(ctx, req)
		is.NoError(err)
		is.NotNil(empty)

		q := "some"

		query := &proto.SearchQuery{
			Query: q,
		}

		res, err := srv.Query(ctx, query)
		is.NoError(err)
		is.Len(res.Documents, 1)
		is.Equal(res.Documents[0].Id, doc.Id)
	})

	t.Run("Shutdown", func(t *testing.T) {
		srv.ShutDown()
	})
}
