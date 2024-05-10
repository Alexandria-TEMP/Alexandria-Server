package models

import (
	"mime/multipart"
)

type Repository struct {
	QuartoProject multipart.File `swaggerignore:"true"`
}
