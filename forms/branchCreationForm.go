package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchCreationForm struct {
	// TODO New files to add to the version

	// Changes made by the branch
	UpdatedPostTitle           *string                           `json:"updatedPostTitle" example:"Updated Project Post Title"`
	UpdatedCompletionStatus    *models.ProjectCompletionStatus   `json:"updatedCompletionStatus" example:"completed"`
	UpdatedScientificFieldIDs  []uint                            `json:"updatedScientificFieldIDs" example:"1"`
	UpdatedFeedbackPreferences *models.ProjectFeedbackPreference `json:"updatedFeedbackPreferences" example:"formal feedback"`

	// The branch's metadata
	CollaboratingMemberIDs []uint `json:"collaboratingMemberIDs" example:"1"`
	ProjectPostID          uint   `json:"projectPostID" example:"1"`
	BranchTitle            string `json:"branchTitle" example:"Proposed Changes"`
	Anonymous              bool   `json:"anonymous" example:"false"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *BranchCreationForm) IsValid() bool {
	exampleString := ""

	return form.UpdatedCompletionStatus.IsValid() &&
		form.UpdatedFeedbackPreferences.IsValid() &&
		form.UpdatedPostTitle != &exampleString &&
		form.BranchTitle != ""
}
