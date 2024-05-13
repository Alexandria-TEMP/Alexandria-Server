package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Connect to the database, auto-migrate its models and return it.
func InitializeDatabase() (*gorm.DB, error) {
	db, err := connectToDatabase()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = autoMigrateAllModels(db)
	if err != nil {
		return nil, fmt.Errorf("could not migrate models: %w", err)
	}

	return db, nil
}
