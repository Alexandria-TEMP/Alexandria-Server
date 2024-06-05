package forms

import (
	"mime/multipart"
)

type IncomingFileForm struct {
	// TODO this doesn't parse right in the swagger thing
	File *multipart.FileHeader `form:"file"`
}
