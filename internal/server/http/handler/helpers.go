package handler

import (
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

	}
}

func getIncr(c *gin.Context) int64 {
	return c.MustGet("incr").(int64)
}

func SetValue(c *gin.Context) {
	valRaw := c.Query("value")
	if valRaw == "" {
		res := gin.H{
			"error": "no value supplied",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	val := []byte(valRaw)

	c.Set("value", val)
}

func getValue(c *gin.Context) []byte {
	return c.MustGet("value").([]byte)
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
