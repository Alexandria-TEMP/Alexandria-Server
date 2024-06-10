package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./postService_interface.go -destination=../../mocks/postService_mock.go

type PostService interface {
	GetPost(postID uint) (*models.Post, error)
	CreatePost(form *forms.PostCreationForm) (*models.Post, error)
	UpdatePost(updatedPost *models.Post) error

	// Return a filtered list of post IDs
	Filter(page, size int, form forms.FilterForm) ([]uint, error)
}
