package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

type PostService struct {
	PostRepository   database.ModelRepositoryInterface[*models.Post]
	MemberRepository database.ModelRepositoryInterface[*models.Member]
	Filesystem       filesystemInterfaces.Filesystem

	PostCollaboratorService interfaces.PostCollaboratorService
}

func (postService *PostService) GetPost(id uint) (*models.Post, error) {
	return postService.PostRepository.GetByID(id)
}

func (postService *PostService) CreatePost(form *forms.PostCreationForm) (*models.Post, error) {
	// Posts created via this function may not be project posts
	// (those must use ProjectPostCreationForms)
	if form.PostType == models.Project {
		return nil, fmt.Errorf("creating post of type ProjectPost using CreatePost is forbidden")
	}

	postCollaborators, err := postService.PostCollaboratorService.MembersToPostCollaborators(form.AuthorMemberIDs, form.Anonymous, models.Author)
	if err != nil {
		return nil, fmt.Errorf("could not create post: %w", err)
	}

	post := models.Post{
		Collaborators:    postCollaborators,
		Title:            form.Title,
		PostType:         form.PostType,
		ScientificFields: form.ScientificFields,
		DiscussionContainer: models.DiscussionContainer{
			// The discussion list is initially empty
			Discussions: []*models.Discussion{},
		},
	}

	if err := postService.PostRepository.Create(&post); err != nil {
		return nil, fmt.Errorf("could not create post: %w", err)
	}

	// TODO filesystem: checkout directory

	// TODO filesystem: create repository

	return &post, nil
}

func (postService *PostService) UpdatePost(_ *models.Post) error {
	// TODO: Access repo to update post here
	return nil
}

/*
	Uploading Post (not ProjectPost) content:
	- Requires having PostID
	- CheckoutDirectory
	- CheckoutBranch("master")
	- SaveZipFile
	- Response 200
	- Start goroutine for rendering
*/
