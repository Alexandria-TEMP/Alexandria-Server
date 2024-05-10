package models

import (
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model

	// Version has one Repository
	VersionID uint

	// TODO
	// QuartoProject multipart.File
}
