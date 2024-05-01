package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/example/hello/reverse"
)

func main() {
	fmt.Println(reverse.String("Hello World!"))
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
