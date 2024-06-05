package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type PostService struct {
	// repo interface here
}

func (postService *PostService) GetPost(_ uint) (*models.Post, error) {
	// TODO: Access repo to get post
	return new(models.Post), nil
}

func (postService *PostService) CreatePost(_ *forms.PostCreationForm) (*models.Post, error) {
	post := &models.Post{
		// TODO fill fields
	}

	// TODO: Add post to repo here

	return post, nil
}

func (postService *PostService) UpdatePost(_ *models.Post) error {
	// TODO: Access repo to update post here
	return nil
}

func (postService *PostService) GetProjectPost(_ uint) (*models.ProjectPost, error) {
	// TODO: Access repo to get post
	return new(models.ProjectPost), nil
}

func (postService *PostService) CreateProjectPost(_ *forms.ProjectPostCreationForm) (*models.ProjectPost, error) {
	post := &models.ProjectPost{
		// TODO fill fields
	}

	// TODO: Add post to repo here

	return post, nil
}

func (postService *PostService) UpdateProjectPost(_ *models.ProjectPost) error {
	// TODO: Access repo to update post here
	return nil
}
