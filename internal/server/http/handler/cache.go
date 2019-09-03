package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucperkins/strato"
)

func (h *Handler) CacheGet(c *gin.Context) {
	log := h.logger("cache/get")

	key := c.Param("key")

	val, err := h.b.CacheGet(key)
	if err != nil {
		if strato.IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.WithError(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	res := gin.H{
		"value": val,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CachePut(c *gin.Context) {
	log := h.logger("cache/put")

	key, ttl := c.Param("key"), getTtl(c)

	value := c.Query("value")
	if value == "" {
		err := gin.H{
			"error": "no cache value provided",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err := h.b.CacheSet(key, value, ttl); err != nil {
		if strato.IsNotFound(err) {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.WithError(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusNoContent)
}
