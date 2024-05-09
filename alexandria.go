package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := router.Run(":8080") // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
}
