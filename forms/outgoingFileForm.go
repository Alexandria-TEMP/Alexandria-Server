package forms

import "io"

type OutgoingFileForm struct {
	File *io.Reader `form:"file"`
}
