package http

import "github.com/gin-gonic/gin"

func (s *Server) routes() *gin.Engine {
	r := gin.New()

	cache := r.Group("/cache/:key")
	{
		cache.GET("", s.cacheGet)

		withTtl := cache.Group("")
		{
			withTtl.Use(setTtl)
			withTtl.PUT("", s.cachePut)
		}
	}

	counters := r.Group("/counters/:counter")
	{
		counters.GET("", s.countersGet)
		counters.PUT("", s.countersPut)
	}

	kv := r.Group("/kv/:key")
	{
		kv.GET("", s.kvGet)
		kv.PUT("", s.kvPut)
		kv.DELETE("", s.kvDelete)
	}

	sets := r.Group("/sets/:set")
	{
		sets.GET("", s.setsGet)

		withItem := sets.Group("")
		{
			withItem.Use(setItem)
			withItem.PUT("", s.setsPut)
			withItem.DELETE("", s.setsDelete)
		}
	}

	return r
}
