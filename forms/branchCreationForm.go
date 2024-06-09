package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type BranchCreationForm struct {
	// TODO New files to add to the version

	// Changes made by the branch
	UpdatedPostTitle           string
	UpdatedCompletionStatus    models.ProjectCompletionStatus
	UpdatedScientificFields    []tags.ScientificField
	UpdatedFeedbackPreferences models.ProjectFeedbackPreference

	// The branch's metadata
	CollaboratingMemberIDs []uint // Get converted to BranchCollaborators
	ProjectPostID          uint
	BranchTitle            string
	Anonymous              bool
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *BranchCreationForm) IsValid() bool {
	return form.UpdatedCompletionStatus.IsValid() &&
		form.UpdatedFeedbackPreferences.IsValid() &&
		len(form.UpdatedPostTitle) >= 0 &&
		len(form.BranchTitle) >= 0
}
