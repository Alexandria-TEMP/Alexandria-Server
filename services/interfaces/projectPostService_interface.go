package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./projectPostService_interface.go -destination=../../mocks/projectPostService_mock.go

type ProjectPostService interface {
	GetProjectPost(postID uint) (*models.ProjectPost, error)
	CreateProjectPost(form *forms.ProjectPostCreationForm, member *models.Member) (*models.ProjectPost, error, error)

	// GetBranchesGroupedByReviewStatus returns branch IDs grouped by their branch review status
	GetBranchesGroupedByReviewStatus(projectPostID uint) (*models.BranchesGroupedByReviewStatusDTO, error)

	// GetDiscussionContainersFromMergeHistory returns discussion containers from the current project version + all previous merged versions
	GetDiscussionContainersFromMergeHistory(projectPostID uint) (*models.DiscussionContainerProjectHistoryDTO, error)
}
