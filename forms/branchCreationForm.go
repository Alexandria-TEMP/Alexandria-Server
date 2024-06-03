package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type BranchCreationForm struct {
	NewPostTitle string

	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields []tags.ScientificField

	Collaborators []*models.BranchCollaborator

	ProjectPostID uint

	BranchTitle string

	Anonymous bool
}
