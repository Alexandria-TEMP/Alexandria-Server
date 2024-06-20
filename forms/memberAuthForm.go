package forms

type MemberAuthForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *MemberAuthForm) IsValid() bool {
	return form.Email != "" &&
		form.Password != ""
}
