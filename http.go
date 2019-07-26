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

	kv := r.Group("/kv")
	{
		kv.GET("/:key", s.kvGet)
		kv.PUT("/:key/:value", s.kvPut)
		kv.DELETE("/:key", s.kvDelete)
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
		if IsNoCacheValue(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
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

func (s *HttpServer) kvGet(c *gin.Context) {
	key := c.Param("key")

	loc := &Location{
		Key: key,
	}

	val, err := s.mem.KVGet(loc)
	if err != nil {
		if IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	res := struct {
		Value string `json:"value"`
	}{
		string(val.Content),
	}

	c.JSON(http.StatusOK, res)
}

func (s *HttpServer) kvPut(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("value")

	loc := &Location{
		Key: key,
	}

	val := &Value{
		Content: []byte(value),
	}

	s.mem.KVPut(loc, val)

	c.Header("Location", fmt.Sprintf("/kv/%s", key))
	c.Status(http.StatusCreated)
}

func (s *HttpServer) kvDelete(c *gin.Context) {
	key := c.Param("key")
	loc := &Location{
		Key: key,
	}

	s.mem.KVDelete(loc)

	c.Status(http.StatusAccepted)
}