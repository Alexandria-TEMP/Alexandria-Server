package models

import (
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	// TODO
	// QuartoProject multipart.File
}
