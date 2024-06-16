package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

type ProjectPostCreationForm struct {
	AuthorMemberIDs           []uint                           `json:"authorMemberIDs" example:"1"`
	Title                     string                           `json:"title" example:"Post Title"`
	Anonymous                 bool                             `json:"anonymous" example:"false"`
	ScientificFieldTagIDs     []uint                           `json:"scientificFieldTagIDs" example:"1"`
	ProjectCompletionStatus   models.ProjectCompletionStatus   `json:"projectCompletionStatus" example:"ongoing"`
	ProjectFeedbackPreference models.ProjectFeedbackPreference `json:"projectFeedbackPreference" example:"formal feedback"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *ProjectPostCreationForm) IsValid() bool {
	return form.ProjectCompletionStatus.IsValid() &&
		form.ProjectFeedbackPreference.IsValid() &&
		form.Title != ""
}
