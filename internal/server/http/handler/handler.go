package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucperkins/strato/internal/backend"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	b   backend.Interface
	log *logrus.Entry
}

func NewHandler(backend backend.Interface, log *logrus.Entry) *Handler {
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
