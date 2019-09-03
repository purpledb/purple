package http

import (
	"fmt"
	"net/http"

	"github.com/lucperkins/strato/internal/server/http/handler"

	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// The core struct undergirding the Strato HTTP interface.
type Server struct {
	addr string
	h    *handler.Handler
	log  *logrus.Entry
}

// Instantiates a new Strato HTTP server using the supplied ServerConfig object.
func NewServer(cfg *strato.ServerConfig) (*Server, error) {
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

func getLogger(cfg *strato.ServerConfig) *logrus.Entry {
	log := logrus.New()

	if cfg.Debug {
		log.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return log.WithField("server", "http")
}

// Starts the Strato HTTP server on the specified port.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.routes(),
	}

	s.log.Infof("starting the Strato HTTP server on %s", s.addr)

	return srv.ListenAndServe()
}
