package database

import (
	"fmt"
	"os"
)

// Information required to connect to the database.
type ConnectionInfo struct {
	username     string
	password     string
	host         string
	databaseName string
	port         int
}

// Read database credentials from environment variables.
// This is only appropriate for a prototype, and not at all safe for production!
func readDatabaseCredentials() (*ConnectionInfo, error) {
	user, err := tryReadingEnv("MARIADB_USER")
	if err != nil {
		return nil, err
	}

	password, err := tryReadingEnv("MARIADB_PASSWORD")
	if err != nil {
		return nil, err
	}

	databaseName, err := tryReadingEnv("MARIADB_DATABASE")
	if err != nil {
		return nil, err
	}

	host, err := tryReadingEnv("ALEXANDRIA_DB_HOST")
	if err != nil {
		return nil, err
	}

	// The port is hard-coded, as it's the standard value to use.
	port := 3306

	info := ConnectionInfo{
		user,
		password,
		host,
		databaseName,
		port,
	}

	return &info, nil
}

// Reads credentials for connecting to the testing database.
func readTestDatabaseCredentials() (*ConnectionInfo, error) {
	// This doesn't re-use the regular database credential fetching code, because
	// these environment variables may be changed down the line.
	user, err := tryReadingEnv("MARIADB_USER")
	if err != nil {
		return nil, err
	}

	password, err := tryReadingEnv("MARIADB_PASSWORD")
	if err != nil {
		return nil, err
	}

	databaseName, err := tryReadingEnv("ALEXANDRIA_TEST_DB_NAME")
	if err != nil {
		return nil, err
	}

	host, err := tryReadingEnv("ALEXANDRIA_TEST_DB_HOST")
	if err != nil {
		return nil, err
	}

	port := 3306

	info := ConnectionInfo{
		user,
		password,
		host,
		databaseName,
		port,
	}

	return &info, nil
}

// Attempt to get an environment variable, return a descriptive error otherwise.
func tryReadingEnv(variable string) (string, error) {
	value, found := os.LookupEnv(variable)

	if !found {
		var zero string
		return zero, fmt.Errorf("environment variable '%s' not present in environment", variable)
	}

	return value, nil
}
