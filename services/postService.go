package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

type PostService struct {
	PostRepository   database.ModelRepositoryInterface[*models.Post]
	MemberRepository database.ModelRepositoryInterface[*models.Member]
	Filesystem       filesystemInterfaces.Filesystem

	PostCollaboratorService interfaces.PostCollaboratorService
	RenderService           interfaces.RenderService
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
		RenderStatus: models.Success,
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

	// clean directory to remove all files
	if err := postService.Filesystem.CleanDir(); err != nil {
		return err
	}

	// save zipped project
	if err := postService.Filesystem.SaveZipFile(c, file); err != nil {
		// it fails so we set render status to failed and reset the branch
		post.RenderStatus = models.Failure
		_, _ = postService.PostRepository.Update(post)
		_ = postService.Filesystem.Reset()

		return fmt.Errorf("failed to save zip file")
	}

	// commit (or perhaps only commit after rendering?)
	if err := postService.Filesystem.CreateCommit(); err != nil {
		return err
	}

	// Set render status pending
	post.RenderStatus = models.Pending
	if _, err := postService.PostRepository.Update(post); err != nil {
		return fmt.Errorf("failed to update post entity")
	}

	go postService.RenderService.RenderPost(post)

	return nil
}

func (postService *PostService) GetMainProject(postID uint) (string, error) {
	var filePath string

	// check post exists
	_, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return filePath, fmt.Errorf("failed to find post with id %v", postID)
	}

	// select repository of the parent post
	postService.Filesystem.CheckoutDirectory(postID)

	// checkout specified branch
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		return filePath, fmt.Errorf("failed to find master branch")
	}

	return postService.Filesystem.GetCurrentZipFilePath(), nil
}

func (postService *PostService) GetMainFiletree(postID uint) (map[string]int64, error, error) {
	// check post exists
	_, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return nil, fmt.Errorf("failed to find post with id %v", postID), nil
	}

	// select repository of the parent post
	postService.Filesystem.CheckoutDirectory(postID)

	// checkout specified branch
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		return nil, fmt.Errorf("failed to find master branch"), nil
	}

	// get file tree
	fileTree, err := postService.Filesystem.GetFileTree()

	return fileTree, nil, err
}

func (postService *PostService) GetMainFileFromProject(postID uint, relFilepath string) (string, error) {
	var absFilepath string

	// validate file path is inside of repository
	if strings.Contains(relFilepath, "..") {
		return absFilepath, fmt.Errorf("file is outside of repository")
	}

	// check post exists
	_, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return absFilepath, fmt.Errorf("failed to find post with id %v", postID)
	}

	// select repository of the post
	postService.Filesystem.CheckoutDirectory(postID)

	// checkout master
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		return absFilepath, fmt.Errorf("failed to find master branch")
	}

	absFilepath = filepath.Join(postService.Filesystem.GetCurrentQuartoDirPath(), relFilepath)

	// Check that file exists, if not return 404
	if exists := utils.FileExists(absFilepath); !exists {
		return "", fmt.Errorf("no such file exists")
	}

	return absFilepath, nil
}

func (postService *PostService) Filter(page, size int, _ forms.FilterForm) ([]uint, error) {
	// TODO construct query based off filter form
	// Future changes: make sure to exclude any posts of type 'Project' from the result!
	// Posts are composed into Project Posts, and those composed Posts shouldn't be returned.
	posts, err := postService.PostRepository.QueryPaginated(page, size, "post_type != 'project'")
	if err != nil {
		return nil, err
	}

	// Extract IDs from the list of posts
	ids := make([]uint, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}

	return ids, nil
}
