package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./postService_interface.go -destination=../../mocks/postService_mock.go

type PostService interface {
	GetPost(postID uint64) (*models.Post, error)
	CreatePost(form *forms.PostCreationForm) *models.Post
	UpdatePost(updatedPost *models.Post) error

	GetProjectPost(postID uint64) (*models.ProjectPost, error)
	CreateProjectPost(form *forms.ProjectPostCreationForm) *models.ProjectPost
	UpdateProjectPost(updatedPost *models.ProjectPost) error
}
