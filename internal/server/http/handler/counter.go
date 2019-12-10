package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CounterGet(c *gin.Context) {
	log := h.logger("counter/get")

	key := c.Param("key")

	val, err := h.b.CounterGet(key)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	res := gin.H{
		"key":   key,
		"value": val,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CounterPut(c *gin.Context) {
	log := h.logger("counter/put")

	key, incr := c.Param("key"), getIncr(c)

	count, err := h.b.CounterIncrement(key, incr)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"counter": key,
		"value": count,
	})
}
