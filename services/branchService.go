package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchService struct {
	BranchRepository database.RepositoryInterface[*models.Branch]
	PostRepository   database.RepositoryInterface[*models.Post]
	Filesystem       filesystemInterfaces.Filesystem
}

func (branchService *BranchService) CreateRepository(postID uint) (models.Branch, error, error) {
	branch := models.Branch{}

	// verify post exists
	if _, err := branchService.PostRepository.GetByID(postID); err != nil {
		return branch, fmt.Errorf("no such post exists"), nil
	}

	// create repository
	if err := branchService.Filesystem.CreateRepository(postID); err != nil {
		return branch, nil, fmt.Errorf("failed to create repository")
	}

	// create branch in db
	if err := branchService.BranchRepository.Create(&branch); err != nil {
		return branch, nil, fmt.Errorf("failed to add branch to db")
	}

	return branch, nil, nil
}

func (branchService *BranchService) CreateBranch(branchCreationForm forms.BranchCreationForm) (models.Branch, error, error) {
	// create  branch
	branch := models.Branch{
		NewPostTitle:            branchCreationForm.NewPostTitle,
		UpdatedCompletionStatus: branchCreationForm.UpdatedCompletionStatus,
		UpdatedScientificFields: branchCreationForm.UpdatedScientificFields,
		Collaborators:           branchCreationForm.Collaborators,
		PostID:                  branchCreationForm.PostID,
		FromBranchID:            branchCreationForm.FromBranchId,
		BranchTitle:             branchCreationForm.BranchTitle,
		Anonymous:               branchCreationForm.Anonymous,
	}

	// verify fromBranch exists
	if _, err := branchService.BranchRepository.GetByID(branchCreationForm.FromBranchId); err != nil {
		return branch, fmt.Errorf("no such fromBranch exists"), nil
	}

	// verify post exists
	if _, err := branchService.PostRepository.GetByID(branchCreationForm.PostID); err != nil {
		return branch, fmt.Errorf("no such post exists"), nil
	}

	// create new git branch in the vfs
	if err := branchService.Filesystem.CheckoutRepository(branch.PostID); err != nil {
		return branch, nil, fmt.Errorf("failed to checkout repository")
	}

	if err := branchService.Filesystem.CreateBranch(branch.FromBranchID); err != nil {
		return branch, nil, fmt.Errorf("failed to create repository")
	}

	// save branch

	if err := branchService.BranchRepository.Create(&branch); err != nil {
		return branch, nil, fmt.Errorf("failed to add branch to db")
	}

	return branch, nil, nil
}
