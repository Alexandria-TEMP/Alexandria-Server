package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type PostCreationForm struct {
	AuthorMemberIDs       []uint          `json:"authorMemberIDs" example:"1"`
	Title                 string          `json:"title" example:"Post Title"`
	Anonymous             bool            `json:"anonymous" example:"false"`
	PostType              models.PostType `json:"postType" example:"question"`
	ScientificFieldTagIDs []uint          `json:"scientificFieldTagIDs" example:"1"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *PostCreationForm) IsValid() bool {
	return form.Title != "" && form.PostType.IsValid()
}
