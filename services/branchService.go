package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchService struct {
	BranchRepository      database.RepositoryInterface[*models.Branch]
	ProjectPostRepository database.RepositoryInterface[*models.ProjectPost]
	Filesystem            filesystemInterfaces.Filesystem
}

func (branchService *BranchService) CreateBranch(branchCreationForm forms.BranchCreationForm) (models.Branch, error, error) {
	var branch models.Branch

	// verify fromBranch exists
	projectPost, err := branchService.ProjectPostRepository.GetByID(branchCreationForm.ProjectPostID)

	if err != nil {
		return branch, fmt.Errorf("no such project post exists"), nil
	}

	branch = models.Branch{
		NewPostTitle:            branchCreationForm.NewPostTitle,
		UpdatedCompletionStatus: branchCreationForm.UpdatedCompletionStatus,
		UpdatedScientificFields: branchCreationForm.UpdatedScientificFields,
		Collaborators:           branchCreationForm.Collaborators,
		ProjectPost:             *projectPost,
		BranchTitle:             branchCreationForm.BranchTitle,
		Anonymous:               branchCreationForm.Anonymous,
	}

	// set vfs to repsitory according to the Post of the ProjectPost of the Branch entity
	branchService.Filesystem.CheckoutDirectory(branch.ProjectPost.PostID)

	// create new branch in git repo with branch ID as its name
	if err := branchService.Filesystem.CreateBranch(string(branch.ID)); err != nil {
		return branch, nil, err
	}

	// save branch entity
	if err := branchService.BranchRepository.Create(&branch); err != nil {
		return branch, nil, fmt.Errorf("failed to add branch to db")
	}

	return branch, nil, nil
}
