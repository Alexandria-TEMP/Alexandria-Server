package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

//go:generate mockgen -source=./postService_interface.go -destination=../../mocks/postService_mock.go

type PostService interface {
	GetPost(postID uint64) (*models.Post, error)
	CreatePost(form *forms.PostCreationForm) *models.Post
	UpdatePost(updatedPost *models.Post) error

	GetProjectPost(postID uint64) (*models.ProjectPost, error)
	CreateProjectPost(form *forms.ProjectPostCreationForm) *models.ProjectPost
	UpdateProjectPost(updatedPost *models.ProjectPost) error
}
