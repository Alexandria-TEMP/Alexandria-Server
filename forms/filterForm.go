package forms

type FilterForm struct {
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *FilterForm) IsValid() bool {
	return true
}
