package forms

import (
	"mime/multipart"
)

type IncomingFileForm struct {
	File *multipart.FileHeader `form:"file"`
}
