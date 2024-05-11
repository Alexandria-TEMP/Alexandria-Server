package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

type PostService struct {
	// repo interface here
}

func (postService *PostService) GetPost(_ uint64) *models.Post {
	// TODO: Access repo to get post
	return new(models.Post)
}

func (postService *PostService) CreatePost(form *forms.PostCreationForm) *models.Post {
	post := &models.Post{
		ID:             0, // Generate uuid here
		PostMetadata:   form.PostMetadata,
		CurrentVersion: form.CurrentVersion,
	}

	// TODO: Add post to repo here

	return post
}

func (postService *PostService) GetProjectPost(_ uint64) *models.ProjectPost {
	// TODO: Access repo to get post
	return new(models.ProjectPost)
}

func (postService *PostService) CreateProjectPost(form *forms.ProjectPostCreationForm) *models.ProjectPost {
	post := &models.ProjectPost{
		ID:                  0, // Generate uuid herePost
		ProjectMetadata:     form.ProjectMetadata,
		OpenMergeRequests:   form.OpenMergeRequests,
		ClosedMergeRequests: form.ClosedMergeRequests,
	}

	// TODO: Add post to repo here

	return post
}
