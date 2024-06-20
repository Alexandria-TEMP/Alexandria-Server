package forms

type DiscussionCreationForm struct {
	// If anonymous, the discussion will ignore member ID
	Anonymous bool   `json:"anonymous" example:"false"`
	Text      string `json:"text" example:"Discussion content."`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *DiscussionCreationForm) IsValid() bool {
	return form.Text != ""
}

type RootDiscussionCreationForm struct {
	// The DiscussionContainer this Discussion will be added to
	ContainerID uint `json:"containerID" example:"1"`

	DiscussionCreationForm DiscussionCreationForm `json:"discussion"`
}

func (form *RootDiscussionCreationForm) IsValid() bool {
	return form.DiscussionCreationForm.IsValid()
}

type ReplyDiscussionCreationForm struct {
	// The Discussion this Discussion will be added to
	ParentID uint `json:"parentID" example:"1"`

	DiscussionCreationForm DiscussionCreationForm `json:"discussion"`
}

func (form *ReplyDiscussionCreationForm) IsValid() bool {
	return form.DiscussionCreationForm.IsValid()
}
