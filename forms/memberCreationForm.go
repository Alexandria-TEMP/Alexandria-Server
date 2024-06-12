package forms

type MemberCreationForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	// making the password just a string for now
	// TODO: some hashing or semblance of security
	Password              string `json:"password"`
	Institution           string `json:"institution"`
	ScientificFieldTagIDs []uint `json:"scientificFieldTagIDs"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *MemberCreationForm) IsValid() bool {
	return form.FirstName != "" &&
		form.LastName != "" &&
		form.Email != "" && // TODO proper email validation
		form.Password != "" &&
		form.Institution != ""
}
