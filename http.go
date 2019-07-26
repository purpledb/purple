package strato

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HttpServer struct {
	mem *Memory
}

func NewHttpServer() *HttpServer {
	mem := NewMemory()

	return &HttpServer{
		mem: mem,
	}
}

func (s *HttpServer) Start() error {
	srv := &http.Server{
		Addr: ":8081",
		Handler: s.routes(),
	}

	return srv.ListenAndServe()
}

func (s *HttpServer) routes() *gin.Engine {
	r := gin.New()

	cache := r.Group("/cache")
	{
		cache.GET("", s.cacheGet)
		cache.PUT("", s.cachePut)
	}

	return r
}

func (s *HttpServer) cacheGet(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.String(http.StatusBadRequest, "no key provided")
		return
	}

	val, err := s.mem.CacheGet(key)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Value string `json:"value"`
	}{
		val,
	}

	c.JSON(http.StatusOK, res)
}

func (s *HttpServer) cachePut(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.String(http.StatusBadRequest, "no key provided")
		return
	}

	value := c.Query("value")
	if value == "" {
		c.String(http.StatusBadRequest, "no value provided")
		return
	}

	ttlStr := c.Query("ttl")
	if ttlStr == "" {
		c.String(http.StatusBadRequest, "no TTL provided")
		return
	}

	ttlInt, err := strconv.Atoi(ttlStr)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("could not convert %s to TTL integer", ttlStr))
		return
	}

	ttl := int32(ttlInt)

	if err := s.mem.CacheSet(key, value, ttl); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}