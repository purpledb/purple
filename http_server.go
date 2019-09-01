package strato

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	address string
	mem     *Memory
	log     *logrus.Entry
}

func NewHttpServer(cfg *HttpConfig) *HttpServer {
	addr := fmt.Sprintf(":%d", cfg.Port)

	mem := NewMemoryBackend()

	logger := logrus.New()

	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.WithField("server", "http")

	return &HttpServer{
		address: addr,
		mem:     mem,
		log:     log,
	}
}

func (s *HttpServer) Start() error {
	srv := &http.Server{
		Addr:    s.address,
		Handler: s.routes(),
	}

	s.log.Infof("starting the Strato HTTP server on %s", s.address)

	return srv.ListenAndServe()
}

func (s *HttpServer) routes() *gin.Engine {
	r := gin.New()

	cache := r.Group("/cache")
	{
		cache.GET("", s.cacheGet)
		cache.PUT("", s.cachePut)
	}

	counters := r.Group("/counters")
	{
		counters.GET("/:counter", s.countersGet)
		counters.PUT("/:counter", s.countersPut)
	}

	kv := r.Group("/kv")
	{
		kv.GET("/:key", s.kvGet)
		kv.PUT("/:key/:value", s.kvPut)
		kv.DELETE("/:key", s.kvDelete)
	}

	search := r.Group("/search")
	{
		search.GET("", s.searchGet)
		search.PUT("", s.searchPut)
	}

	sets := r.Group("/sets")
	{
		sets.GET("/:set", s.setsGet)
		sets.PUT("/:set/:item", s.setsPut)
		sets.DELETE("/:set/:item", s.setsDelete)
	}

	return r
}

func (s *HttpServer) cacheGet(c *gin.Context) {
	log := s.log.WithField("op", "cache/get")

	key := c.Query("key")
	if key == "" {
		c.String(http.StatusBadRequest, "no key provided")
		return
	}

	val, err := s.mem.CacheGet(key)
	if err != nil {
		if IsNoItemFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else if IsExpired(err) {
			c.Status(http.StatusGone)
			return
		} else {
			log.Error(err)
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
	log := s.log.WithField("op", "cache/put")

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
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (s *HttpServer) countersGet(c *gin.Context) {
	counter := c.Param("counter")

	value, err := s.mem.CounterGet(counter)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	res := struct {
		Counter string `json:"counter"`
		Value   int32  `json:"value"`
	}{
		Counter: counter,
		Value:   value,
	}

	c.JSON(http.StatusOK, res)
}

func (s *HttpServer) countersPut(c *gin.Context) {
	counter, incr := c.Param("counter"), c.Query("increment")

	if incr == "" {
		c.String(http.StatusBadRequest, "no increment specified")
		return
	}

	i, err := strconv.Atoi(incr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	increment := int32(i)

	if err := s.mem.CounterIncrement(counter, increment); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *HttpServer) kvGet(c *gin.Context) {
	log := s.log.WithField("op", "kv/get")

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
			log.Error(err)
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
	log := s.log.WithField("op", "kv/put")

	key := c.Param("key")
	value := c.Param("value")

	loc := &Location{
		Key: key,
	}

	val := &Value{
		Content: []byte(value),
	}

	if err := s.mem.KVPut(loc, val); err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Location", fmt.Sprintf("/kv/%s", key))
	c.Status(http.StatusCreated)
}

func (s *HttpServer) kvDelete(c *gin.Context) {
	log := s.log.WithField("op", "kv/delete")

	key := c.Param("key")
	loc := &Location{
		Key: key,
	}

	if err := s.mem.KVDelete(loc); err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *HttpServer) searchGet(c *gin.Context) {
	q := c.Query("q")

	if q == "" {
		c.String(http.StatusBadRequest, "no query string provided")
		return
	}

	q = strings.ToLower(q)

	docs := s.mem.Query(q)

	res := struct {
		Query     string      `json:"query"`
		Documents []*Document `json:"documents"`
	}{
		Query:     q,
		Documents: docs,
	}

	c.JSON(http.StatusOK, res)
}

func (s *HttpServer) searchPut(c *gin.Context) {
	var doc Document

	if err := c.ShouldBind(&doc); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	s.mem.Index(&doc)

	c.Status(http.StatusAccepted)
}

func (s *HttpServer) setsGet(c *gin.Context) {
	set := c.Param("set")

	items := s.mem.GetSet(set)

	res := struct {
		Set   string   `json:"set"`
		Items []string `json:"items"`
	}{
		Set:   set,
		Items: items,
	}

	c.JSON(http.StatusOK, res)
}

func (s *HttpServer) setsPut(c *gin.Context) {
	set, item := c.Param("set"), c.Param("item")

	s.mem.AddToSet(set, item)

	c.Status(http.StatusAccepted)
}

func (s *HttpServer) setsDelete(c *gin.Context) {
	set, item := c.Param("set"), c.Param("item")

	s.mem.RemoveFromSet(set, item)

	c.Status(http.StatusAccepted)
}
