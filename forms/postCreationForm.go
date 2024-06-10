package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type PostCreationForm struct {
	// TODO send files somehow?

	AuthorMemberIDs  []uint                   `json:"authorMemberIDs"`
	Title            string                   `json:"title"`
	Anonymous        bool                     `json:"anonymous"`
	PostType         models.PostType          `json:"postType"`
	ScientificFields []models.ScientificField `json:"ScientificFields"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *PostCreationForm) IsValid() bool {
	return form.Title != "" && form.PostType.IsValid()
}
