package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

type ProjectPostCreationForm struct {
	PostCreationForm          PostCreationForm                 `json:"postCreationForm"`
	ProjectCompletionStatus   models.ProjectCompletionStatus   `json:"projectCompletionStatus"`
	ProjectFeedbackPreference models.ProjectFeedbackPreference `json:"projectFeedbackPreference"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *ProjectPostCreationForm) IsValid() bool {
	return form.PostCreationForm.IsValid() &&
		form.ProjectCompletionStatus.IsValid() &&
		form.ProjectFeedbackPreference.IsValid()
}
