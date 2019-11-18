package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purpledb/purple/internal/backend"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	b   backend.Service
	log *logrus.Entry
}

func NewHandler(backend backend.Service, log *logrus.Entry) *Handler {
	return &Handler{
		b:   backend,
		log: log,
	}
}

func (h *Handler) logger(op string) *logrus.Entry {
	return h.log.WithField("op", op)
}

func (h *Handler) Ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
