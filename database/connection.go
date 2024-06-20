package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectToDatabase(info *ConnectionInfo, silent bool) (*gorm.DB, error) {
	connectionString := getConnectionString(info)

	// Choose a database logger (used by GORM)
	var log logger.Interface
	if silent {
		log = logger.Default.LogMode(logger.Silent)
	} else {
		log = logger.Default
	}

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: log,
	})

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
