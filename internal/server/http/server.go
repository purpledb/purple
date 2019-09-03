package http

import (
	"fmt"
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend"
	"net/http"
	"strconv"

	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// The core struct undergirding the Strato HTTP interface.
type Server struct {
	address string
	backend backend.Interface
	log     *logrus.Entry
}

// Instantiates a new Strato HTTP server using the supplied ServerConfig object.
func NewServer(cfg *strato.ServerConfig) (*Server, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)

	bk, err := backend.NewBackend(cfg)
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

	return &Server{
		address: addr,
		backend: bk,
		log:     log,
	}, nil
}

// Starts the Strato HTTP server on the specified port.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    s.address,
		Handler: s.routes(),
	}

	s.log.Infof("starting the Strato HTTP server on %s", s.address)

	return srv.ListenAndServe()
}

// Cache operations

func (s *Server) cacheGet(c *gin.Context) {
	log := s.log.WithField("op", "cache/get")

	key := c.Param("key")

	val, err := s.backend.CacheGet(key)
	if err != nil {
		if strato.IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else if strato.IsExpired(err) {
			c.Status(http.StatusGone)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
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

func (s *Server) cachePut(c *gin.Context) {
	log := s.log.WithField("op", "cache/put")

	key, ttl := c.Param("key"), getTtl(c)

	value := c.Query("value")
	if value == "" {
		c.String(http.StatusBadRequest, "no value provided")
		return
	}

	if err := s.backend.CacheSet(key, value, int32(ttl)); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

// Counter operations

func (s *Server) countersGet(c *gin.Context) {
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

func (s *Server) countersPut(c *gin.Context) {
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

// KV operations

func (s *Server) kvGet(c *gin.Context) {
	log := s.log.WithField("op", "kv/get")

	key := c.Param("key")

	val, err := s.backend.KVGet(key)
	if err != nil {
		if strato.IsNotFound(err) {
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

func (s *Server) kvPut(c *gin.Context) {
	log := s.log.WithField("op", "kv/put")

	key, value := c.Param("key"), c.Query("value")
	if value == "" {
		c.String(http.StatusBadRequest, "no value specified")
		return
	}

	val := &strato.Value{
		Content: []byte(value),
	}

	if err := s.backend.KVPut(key, val); err != nil {
		if strato.IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Header("Location", fmt.Sprintf("/kv/%s", key))
	c.Status(http.StatusCreated)
}

func (s *Server) kvDelete(c *gin.Context) {
	log := s.log.WithField("op", "kv/delete")

	key := c.Param("key")

	if err := s.backend.KVDelete(key); err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

// Sets
func setItem(c *gin.Context) {
	item := c.Query("item")
	if item == "" {
		c.String(http.StatusBadRequest, "no item provided")
		c.Abort()
		return
	}

	c.Set("item", item)
}

func getItem(c *gin.Context) string {
	return c.MustGet("item").(string)
}

func setTtl(c *gin.Context) {
	ttlRaw := c.Query("ttl")
	if ttlRaw == "" {
		c.String(http.StatusBadRequest, "no TTL provided")
		c.Abort()
		return
	}

	ttl, err := strconv.Atoi(ttlRaw)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	c.Set("ttl", ttl)
}

func getTtl(c *gin.Context) int {
	return c.MustGet("ttl").(int)
}

func (s *Server) setsGet(c *gin.Context) {
	set := c.Param("set")

	items, err := s.backend.GetSet(set)
	if err != nil {
		if err == strato.ErrNoSet {
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

func (s *Server) setsPut(c *gin.Context) {
	set := c.Param("set")

	item := getItem(c)

	if err := s.backend.AddToSet(set, item); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *Server) setsDelete(c *gin.Context) {
	set := c.Param("set")

	item := getItem(c)

	if err := s.backend.RemoveFromSet(set, item); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusAccepted)
}
