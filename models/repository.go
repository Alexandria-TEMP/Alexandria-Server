package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	QuartoProject multipart.File
}
