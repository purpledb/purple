package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purpledb/purple"
)

const todosSet = "todos"

func getTodos(client *purple.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := client.SetGet(todosSet)
		if err != nil {
			if purple.IsNotFound(err) {
				c.JSON(http.StatusOK, gin.H{"todos": []string{}})
				return
			} else {
				log.Println(err)
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	}
}

func createTodo(client *purple.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo := getTodo(c)

		todos, err := client.SetAdd(todosSet, todo)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	}
}

func deleteTodo(client *purple.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo := getTodo(c)

		todos, err := client.SetRemove(todosSet, todo)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	}
}

func setTodo(c *gin.Context) {
	todo := c.Query("todo")
	if todo == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "no todo provided",
		})
		return
	}

	c.Set("todo", todo)
}

func getTodo(c *gin.Context) string {
	return c.MustGet("todo").(string)
}

func main() {
	r := gin.Default()

	client, err := purple.NewGrpcClient(&purple.ClientConfig{
		Address: "purple:8080",
	})

	if err != nil {
		log.Fatal(err)
	}

	todos := r.Group("/todos")
	{
		todos.GET("", getTodos(client))

		withTodo := todos.Group("")
		{
			withTodo.Use(setTodo)

			withTodo.POST("", createTodo(client))
			withTodo.DELETE("", deleteTodo(client))
		}
	}

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
