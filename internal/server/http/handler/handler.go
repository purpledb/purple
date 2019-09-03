package handler

import (
	"github.com/lucperkins/strato/internal/backend"
	"github.com/sirupsen/logrus"
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
