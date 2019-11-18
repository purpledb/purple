package http

import (
	"fmt"
	"net/http"

	"github.com/purpledb/purple"

	"github.com/purpledb/purple/internal/server/http/handler"

	"github.com/purpledb/purple/internal/backend"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// The core struct undergirding the purple HTTP interface.
type Server struct {
	addr string
	h    *handler.Handler
	log  *logrus.Entry
}

// Instantiates a new purple HTTP server using the supplied ServerConfig object.
func NewServer(cfg *purple.ServerConfig) (*Server, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)

	bk, err := backend.NewBackend(cfg)
	if err != nil {
		return nil, err
	}

	log := getLogger(cfg)

	h := handler.NewHandler(bk, log)

	return &Server{
		addr: addr,
		h:    h,
		log:  log,
	}, nil
}

func getLogger(cfg *purple.ServerConfig) *logrus.Entry {
	log := logrus.New()

	if cfg.Debug {
		log.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return log.WithField("server", "http")
}

// Starts the purple HTTP server on the specified port.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.routes(),
	}

	s.log.Infof("starting the purple HTTP server on %s", s.addr)

	return srv.ListenAndServe()
}
