package forms

type MemberCreationForm struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"example@example.example"`
	// making the password just a string for now
	// TODO: some hashing or semblance of security
	Password              string `json:"password" example:"password"`
	Institution           string `json:"institution" example:"Example University"`
	ScientificFieldTagIDs []uint `json:"scientificFieldTagIDs" example:"1"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *MemberCreationForm) IsValid() bool {
	return form.FirstName != "" &&
		form.LastName != "" &&
		form.Email != "" && // TODO proper email validation
		form.Password != "" &&
		form.Institution != ""
}
