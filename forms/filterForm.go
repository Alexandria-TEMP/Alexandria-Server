package forms

type PostFilterForm struct {
	IncludeProjectPosts bool `json:"includeProjectPosts" example:"true"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *PostFilterForm) IsValid() bool {
	return true
}

type ProjectPostFilterForm struct {
}

func (form *ProjectPostFilterForm) IsValid() bool {
	return true
}
