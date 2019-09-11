package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lucperkins/strato"
	"log"
	"net/http"
)

func getTodos(client *strato.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := client.SetGet("todos")
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

func createTodo(client *strato.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo := c.Query("todo")
		if todo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no todo provided",
			})
			return
		}

		todos, err := client.SetAdd("todos", todo)
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

func deleteTodo(client *strato.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo := c.Query("todo")
		if todo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no todo provided",
			})
			return
		}

		todos, err := client.SetRemove("todos", todo)
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

func main() {
	r := gin.Default()

	client, err := strato.NewGrpcClient(&strato.ClientConfig{
		Address: "strato:8080",
	})

	if err != nil {
		log.Fatal(err)
	}

	todos := r.Group("/todos")
	{
		todos.POST("", createTodo(client))
		todos.GET("", getTodos(client))
		todos.DELETE("", deleteTodo(client))
	}

	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
