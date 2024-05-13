package models

import (
	"mime/multipart"
)

type Repository struct {
	File *multipart.FileHeader `form:"file"`

	// TODO write serialization/deserialization, OR use a filesystem instead
	// QuartoProject multipart.File `swaggerignore:"true"`
}
