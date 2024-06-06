package forms

type MemberCreationForm struct {
	FirstName string
	LastName  string
	Email     string
	// making the password just a string for now
	// TODO: some hashing or semblance of security
	Password              string
	Institution           string
	ScientificFieldTagIDs []uint
}
