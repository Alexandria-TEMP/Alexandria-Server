package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectToDatabase(info *ConnectionInfo) (*gorm.DB, error) {
	connectionString := getConnectionString(info)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// Convert database connection info into a connection string.
func getConnectionString(info *ConnectionInfo) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		info.username, info.password, info.host, info.port, info.databaseName)
}
