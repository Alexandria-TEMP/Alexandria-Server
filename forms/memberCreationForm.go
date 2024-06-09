package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"

type MemberCreationForm struct {
	FirstName string
	LastName  string
	Email     string
	// making the password just a string for now
	// TODO: some hashing or semblance of security
	Password    string
	Institution string
	Fields      []tags.ScientificFieldTag
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *MemberCreationForm) IsValid() bool {
	return form.FirstName != "" &&
		form.LastName != "" &&
		form.Email != "" && // TODO proper email validation
		form.Password != "" &&
		form.Institution != ""
}
