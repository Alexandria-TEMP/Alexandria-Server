package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Connect to the database, auto-migrate its models and return it.
func InitializeDatabase() (*gorm.DB, error) {
	info, err := readDatabaseCredentials()
	if err != nil {
		return nil, fmt.Errorf("could not read database credentials: %w", err)
	}

	db, err := connectToDatabase(info, false)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = autoMigrateAllModels(db)
	if err != nil {
		return nil, fmt.Errorf("could not migrate models: %w", err)
	}

	return db, nil
}

func InitializeTestDatabase() (*gorm.DB, error) {
	// Since the test database should absolutely never be confused for the production
	// database, this function is entirely separated and doesn't re-use the above code.
	info, err := readTestDatabaseCredentials()
	if err != nil {
		return nil, fmt.Errorf("could not read database credentials: %w", err)
	}

	db, err := connectToDatabase(info, true)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = autoMigrateAllModels(db)
	if err != nil {
		return nil, fmt.Errorf("could not migrate models: %w", err)
	}

	return db, nil
}
