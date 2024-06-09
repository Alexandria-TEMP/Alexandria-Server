package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./projectPostService_interface.go -destination=../../mocks/projectPostService_mock.go

type ProjectPostService interface {
	GetProjectPost(postID uint) (*models.ProjectPost, error)
	CreateProjectPost(form *forms.ProjectPostCreationForm) (*models.ProjectPost, error)
	UpdateProjectPost(updatedPost *models.ProjectPost) error
}
