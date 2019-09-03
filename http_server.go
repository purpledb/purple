package strato

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	address string
	backend Backend
	log     *logrus.Entry
}

func NewHttpServer(cfg *ServerConfig) (*HttpServer, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)

	backend, err := NewBackend(cfg)
	if err != nil {
		return nil, err
	}

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
		backend: backend,
		log:     log,
	}, nil
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
		kv.GET("/:bucket/:key", s.kvGet)
		kv.PUT("/:bucket/:key/:value", s.kvPut)
		kv.DELETE("/:bucket/:key", s.kvDelete)
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

	val, err := s.backend.CacheGet(key)
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

	if err := s.backend.CacheSet(key, value, ttl); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (s *HttpServer) countersGet(c *gin.Context) {
	counter := c.Param("counter")

	value, err := s.backend.CounterGet(counter)
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
		Value   int64  `json:"value"`
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
		return
	}

	if err := s.backend.CounterIncrement(counter, int64(i)); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

// KV
func getLocation(c *gin.Context) *Location {
	return &Location{
		Bucket: c.Param("bucket"),
		Key:    c.Param("key"),
	}
}

func (s *HttpServer) kvGet(c *gin.Context) {
	log := s.log.WithField("op", "kv/get")

	loc := getLocation(c)

	val, err := s.backend.KVGet(loc)
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

	loc := getLocation(c)

	value := c.Param("value")

	val := &Value{
		Content: []byte(value),
	}

	if err := s.backend.KVPut(loc, val); err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Location", fmt.Sprintf("/kv/%s/%s", loc.Bucket, loc.Key))
	c.Status(http.StatusCreated)
}

func (s *HttpServer) kvDelete(c *gin.Context) {
	log := s.log.WithField("op", "kv/delete")

	loc := getLocation(c)

	if err := s.backend.KVDelete(loc); err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *HttpServer) setsGet(c *gin.Context) {
	set := c.Param("set")

	items, err := s.backend.GetSet(set)
	if err != nil {
		if err == ErrNoSet {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

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

	if err := s.backend.AddToSet(set, item); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *HttpServer) setsDelete(c *gin.Context) {
	set, item := c.Param("set"), c.Param("item")

	if err := s.backend.RemoveFromSet(set, item); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}
