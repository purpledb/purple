package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucperkins/strato"
)

func (h *Handler) SetGet(c *gin.Context) {
	log := h.logger("set/get")

	key := c.Param("key")

	items, err := h.b.GetSet(key)
	if err != nil {
		if err == strato.ErrNoSet {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	res := gin.H{
		"set":   key,
		"items": items,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) SetPut(c *gin.Context) {
	log := h.logger("set/increment")

	key, item := c.Param("key"), getItem(c)

	items, err := h.b.AddToSet(key, item)
	if err != nil {
		if err == strato.ErrNoSet {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	res := gin.H{
		"set":   key,
		"items": items,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) SetDelete(c *gin.Context) {
	log := h.logger("set/remove")

	key, item := c.Param("key"), getItem(c)

	items, err := h.b.RemoveFromSet(key, item)
	if err != nil {
		if err == strato.ErrNoSet {
			c.Status(http.StatusNotFound)
			return
		} else {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	res := gin.H{
		"set":   key,
		"items": items,
	}

	c.JSON(http.StatusOK, res)
}
