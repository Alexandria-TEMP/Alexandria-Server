package database

import (
	"fmt"
)

func InitializeDatabase() error {
	// If DB connection fails, terminate
	db, err := database.ConnectToDatabase()
	if err != nil {
		return fmt.Print("could not connect to database: %s", err)
	}

	err = database.AutoMigrateAllModels(db)
	if err != nil {
		log.Fatalf("could not migrate models: %s", err)
	}
}
