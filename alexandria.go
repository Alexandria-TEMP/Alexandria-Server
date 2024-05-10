package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
)

func main() {
	// If DB connection fails, terminate
	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	err = database.AutoMigrateAllModels(db)
	if err != nil {
		log.Fatalf("could not migrate models: %s", err)
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err = router.Run(":8080") // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
}
