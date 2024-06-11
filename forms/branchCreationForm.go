package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchCreationForm struct {
	// TODO New files to add to the version

	// Changes made by the branch
	UpdatedPostTitle           *string                           `json:"updatedPostTitle"`
	UpdatedCompletionStatus    *models.ProjectCompletionStatus   `json:"updatedCompletionStatus"`
	UpdatedScientificFields    []models.ScientificField          `json:"updatedScientificFields"`
	UpdatedFeedbackPreferences *models.ProjectFeedbackPreference `json:"updatedFeedbackPreferences"`

	// The branch's metadata
	CollaboratingMemberIDs []uint `json:"collaboratingMemberIDs"`
	ProjectPostID          uint   `json:"projectPostID"`
	BranchTitle            string `json:"branchTitle"`
	Anonymous              bool   `json:"anonymous"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *BranchCreationForm) IsValid() bool {
	exampleString := ""

	return form.UpdatedCompletionStatus.IsValid() &&
		form.UpdatedFeedbackPreferences.IsValid() &&
		form.UpdatedPostTitle != &exampleString &&
		form.BranchTitle != ""
}
