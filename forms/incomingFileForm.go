package forms

import (
	"mime/multipart"
)

type IncomingFileForm struct {
	// TODO this doesn't parse right in the swagger thing
	File *multipart.FileHeader `form:"file"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *IncomingFileForm) IsValid() bool {
	return true
}
