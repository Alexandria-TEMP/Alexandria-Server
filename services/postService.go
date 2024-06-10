package services

import (
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
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

	// Checkout directory where post will store it's files
	postService.Filesystem.CheckoutDirectory(post.ID)

	// Create a new git repo there
	if err := postService.Filesystem.CreateRepository(); err != nil {
		return nil, err
	}

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

func (postService *PostService) UploadPost(c *gin.Context, file *multipart.FileHeader, postID uint) error {
	// get post
	post, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return fmt.Errorf("failed to find postID with id %v", postID)
	}

	// select repository of the post and checkout master
	postService.Filesystem.CheckoutDirectory(postID)
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		return err
	}

	return nil
}
