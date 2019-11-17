package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucperkins/purple"
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

	if err := h.b.CounterIncrement(key, incr); err != nil {
		if purple.IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusNoContent)
}
