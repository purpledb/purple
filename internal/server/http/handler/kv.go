package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/services/kv"
)

func (h *Handler) KvGet(c *gin.Context) {
	log := h.logger("kv/get")

	key := c.Param("key")

	val, err := h.b.KVGet(key)
	if err != nil {
		if strato.IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	res := gin.H{
		"value": string(val.Content),
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) KvPut(c *gin.Context) {
	log := h.logger("kv/put")

	key, val := c.Param("key"), getValue(c)

	value := &kv.Value{
		Content: val,
	}

	if err := h.b.KVPut(key, value); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) KvDelete(c *gin.Context) {
	log := h.logger("kv/delete")

	key := c.Param("key")

	if err := h.b.KVDelete(key); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}