package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchService interface {
	// CreateBranch creates a new branch from a creation form.
	// It assumes that a repository has already been created for this post.
	// Error 1 404
	// Error 2 500
	CreateBranch(branchCreationForm forms.BranchCreationForm) (models.Branch, error, error)

	// GetBranch gets an existing branch from the DB
	GetBranch(branchID uint) (models.Branch, error)
}
