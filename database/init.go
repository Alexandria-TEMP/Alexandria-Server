package database

import (
	"fmt"

	"gorm.io/gorm"
)

func InitializeDatabase() (*gorm.DB, error) {
	db, err := ConnectToDatabase()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = AutoMigrateAllModels(db)
	if err != nil {
		return nil, fmt.Errorf("could not migrate models: %w", err)
	}

	return db, nil
}
