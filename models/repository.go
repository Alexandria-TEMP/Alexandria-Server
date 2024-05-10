package models

import (
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model

	// Version has one Repository
	VersionID uint

	// TODO write serialization/deserialization, OR use a filesystem instead
	// QuartoProject multipart.File
}
