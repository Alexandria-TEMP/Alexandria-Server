package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchService interface {
	// CreateBranch creates a new branch from a creation form
	// Error 1 404
	// Error 2 500
	CreateBranch(branchCreationForm forms.BranchCreationForm) (models.Branch, error, error)

	// CreateRepository creates a new repository in the vfs, with a clean main branch
	// Error 1 404
	// Error 2 500
	CreateRepository(postID uint) (models.Branch, error, error)
}
