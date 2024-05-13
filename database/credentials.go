package database

import (
	"fmt"
	"os"
)

// Information required to connect to the database.
type databaseConnectionInfo struct {
	username     string
	password     string
	host         string
	databaseName string
	port         int
}

// Read database credentials from environment variables.
// This is only appropriate for a prototype, and not at all safe for production!
func readDatabaseCredentials() (*databaseConnectionInfo, error) {
	user, found := os.LookupEnv("MARIADB_USER")
	if !found {
		return nil, fmt.Errorf("unable to find environment variable MARIADB_USER")
	}

	password, found := os.LookupEnv("MARIADB_PASSWORD")
	if !found {
		return nil, fmt.Errorf("unable to find environment variable MARIADB_PASSWORD")
	}

	databaseName, found := os.LookupEnv("MARIADB_DATABASE")
	if !found {
		return nil, fmt.Errorf("unable to find environment variable MARIADB_DATABASE")
	}

	host, found := os.LookupEnv("ALEXANDRIA_DB_HOST")
	if !found {
		return nil, fmt.Errorf("unable to find environment variable ALEXANDRIA_DB_HOST")
	}

	// The port is hard-coded, as it's the standard value to use.
	port := 3306

	info := databaseConnectionInfo{
		user,
		password,
		host,
		databaseName,
		port,
	}

	return &info, nil
}
