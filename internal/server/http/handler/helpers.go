package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetTtl(c *gin.Context) {
	ttlRaw := c.Query("ttl")
	if ttlRaw == "" {
		res := gin.H{
			"error": "no TTL provided",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ttl, err := strconv.Atoi(ttlRaw)
	if err != nil {
		res := gin.H{
			"error": err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	c.Set("ttl", int32(ttl))
}

func getTtl(c *gin.Context) int32 {
	return c.MustGet("ttl").(int32)
}

func SetIncr(c *gin.Context) {
	incrRaw := c.Query("increment")
	if incrRaw == "" {
		res := gin.H{
			"error": "no increment specified",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	incr, err := strconv.ParseInt(incrRaw, 10, 64)
	if err != nil {
		res := gin.H{
			"error": fmt.Sprintf("could not parse %s into an integer", incrRaw),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	c.Set("increment", incr)
}

func getIncr(c *gin.Context) int64 {
	return c.MustGet("increment").(int64)
}

type valJs struct {
	Content string `json:"content"`
}

func SetValue(c *gin.Context) {
	var js valJs

	if err := c.ShouldBind(&js); err != nil {
		res := gin.H{
			"error": err.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if js.Content == "" {
		res := gin.H{
			"error": "content cannot be empty",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	c.Set("value", &js)
}

func getValue(c *gin.Context) *valJs {
	return c.MustGet("value").(*valJs)
}

func SetItem(c *gin.Context) {
	item := c.Query("item")
	if item == "" {
		res := gin.H{
			"error": "no item supplied",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	c.Set("item", item)
}

func getItem(c *gin.Context) string {
	return c.MustGet("item").(string)
}
