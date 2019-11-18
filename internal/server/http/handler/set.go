package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purpledb/purple"
)

func emptySetRes(set string) gin.H {
	return gin.H{
		"set":   set,
		"items": []string{},
	}
}

func (h *Handler) SetGet(c *gin.Context) {
	log := h.logger("set/get")

	key := c.Param("key")

	items, err := h.b.SetGet(key)
	if err != nil {
		if purple.IsNotFound(err) {
			c.JSON(http.StatusOK, emptySetRes(key))
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

	items, err := h.b.SetAdd(key, item)
	if err != nil {
		if purple.IsNotFound(err) {
			c.JSON(http.StatusOK, emptySetRes(key))
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

	items, err := h.b.SetRemove(key, item)
	if err != nil {
		if purple.IsNotFound(err) {
			c.JSON(http.StatusOK, emptySetRes(key))
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
