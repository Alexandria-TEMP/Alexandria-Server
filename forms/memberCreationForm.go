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
