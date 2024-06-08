package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type BranchCreationForm struct {
	NewPostTitle string `json:"new_post_title"`

	UpdatedCompletionStatus tags.CompletionStatus  `json:"updated_completion_status"`
	UpdatedScientificFields []tags.ScientificField `json:"updated_scientific_fields"`

	Collaborators []*models.BranchCollaborator `json:"collaborators"`

	ProjectPostID uint `json:"project_post_id"`

	BranchTitle string `json:"branch_title"`

	Anonymous bool `json:"anonymous"`
}
