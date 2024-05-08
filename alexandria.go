package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := r.Run(":8080") // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatalf("impossible to start server: %s", err)
	}
}
