package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"

type BranchCreationForm struct {
	// TODO New files to add to the version

	// Changes made by the MR
	UpdatedPostTitle           string
	UpdatedCompletionStatus    tags.CompletionStatus
	UpdatedScientificFields    []*tags.ScientificFieldTag
	UpdatedFeedbackPreferences tags.FeedbackPreference

	// The MR's metadata
	CollaboratingMemberIDs []uint // Get converted to BranchCollaborators
	ProjectPostID          uint
	BranchTitle            string
	Anonymous              bool
}
