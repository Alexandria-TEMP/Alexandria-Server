package interfaces

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

//go:generate mockgen -package=mocks -source=./postCollaboratorService_interface.go -destination=../../mocks/postCollaboratorService_mock.go

type PostCollaboratorService interface {
	GetPostCollaborator(id uint) (*models.PostCollaborator, error)
	MembersToPostCollaborators(IDs []uint, anonymous bool, collaborationType models.CollaborationType) ([]*models.PostCollaborator, error)
	MergeReviewers(projectPost *models.ProjectPost, reviews []*models.BranchReview) error
	MergeContributors(projectPost *models.ProjectPost, branchCollaborators []*models.BranchCollaborator) error
}
