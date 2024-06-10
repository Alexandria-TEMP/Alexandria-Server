package forms

type DiscussionCreationForm struct {
	// Parent ID / Version ID is sent as an optional query parameter

	// If anonymous, the discussion will ignore member ID
	Anonymous bool `json:"anonymous"`
	MemberID  uint `json:"memberID"`

	Text string `json:"text"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *DiscussionCreationForm) IsValid() bool {
	return form.Text != ""
}
