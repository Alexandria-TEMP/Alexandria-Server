package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/flock"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

type PostService struct {
	PostRepository                        database.ModelRepositoryInterface[*models.Post]
	ProjectPostRepository                 database.ModelRepositoryInterface[*models.ProjectPost]
	MemberRepository                      database.ModelRepositoryInterface[*models.Member]
	ScientificFieldTagContainerRepository database.ModelRepositoryInterface[*models.ScientificFieldTagContainer]
	Filesystem                            filesystemInterfaces.Filesystem

	PostCollaboratorService interfaces.PostCollaboratorService
	RenderService           interfaces.RenderService
	TagService              interfaces.TagService
}

func (postService *PostService) GetPost(id uint) (*models.Post, error) {
	return postService.PostRepository.GetByID(id)
}

func (postService *PostService) CreatePost(form *forms.PostCreationForm, member *models.Member) (*models.Post, error) {
	// Posts created via this function may not be project posts
	// (those must use ProjectPostCreationForms)
	if form.PostType == models.Project {
		return nil, fmt.Errorf("creating post of type ProjectPost using CreatePost is forbidden")
	}

	postCollaborators, err := postService.PostCollaboratorService.MembersToPostCollaborators(form.AuthorMemberIDs, form.Anonymous, models.Author)
	if err != nil {
		return nil, fmt.Errorf("could not create post: %w", err)
	}

	// check if creating member is in authors or post is anonymous
	if !form.Anonymous && !slices.Contains(form.AuthorMemberIDs, member.ID) {
		return nil, fmt.Errorf("the creating member is not in the list of authors. either add the member or set the post to anonymous")
	}

	// convert []uint to []*models.ScientificFieldTag
	tags, err := postService.TagService.GetTagsFromIDs(form.ScientificFieldTagIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get tags from ids: %w", err)
	}

	// create and save the tag container to avoid issues with saving later (preloading stuff?)
	postTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: tags,
	}

	if err := postService.ScientificFieldTagContainerRepository.Create(postTagContainer); err != nil {
		return nil, fmt.Errorf("failed to add tag container to db: %w", err)
	}

	// construct post
	post := models.Post{
		Collaborators:               postCollaborators,
		Title:                       form.Title,
		PostType:                    form.PostType,
		ScientificFieldTagContainer: *postTagContainer,
		DiscussionContainer: models.DiscussionContainer{
			// The discussion list is initially empty
			Discussions: []*models.Discussion{},
		},
		RenderStatus: models.Success, // Render status is success because the default project is prerendered
	}

	if err := postService.PostRepository.Create(&post); err != nil {
		return nil, fmt.Errorf("could not create post: %w", err)
	}

	// lock directory and defer unlocking it
	lock, err := postService.Filesystem.LockDirectory(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock for directory %v: %w", post.ID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// Checkout directory where post will store it's files
	postService.Filesystem.CheckoutDirectory(post.ID)

	// Create a new git repo there
	if err := postService.Filesystem.CreateRepository(); err != nil {
		return nil, err
	}

	return &post, nil
}

func (postService *PostService) UploadPost(c *gin.Context, file *multipart.FileHeader, postID uint) error {
	// get post
	post, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return fmt.Errorf("failed to find postID with id %v: %w", postID, err)
	}

	// reject project posts
	if post.PostType == models.Project {
		return fmt.Errorf("this post is a project post and cannot be directly changed. instead propose a change using a branch")
	}

	// lock directory
	// we unlock it upon error or after finishing render pipeline
	lock, err := postService.Filesystem.LockDirectory(post.ID)
	if err != nil {
		return fmt.Errorf("failed to acquire lock for directory %v: %w", post.ID, err)
	}

	// select repository of the post and checkout master
	postService.Filesystem.CheckoutDirectory(postID)

	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return err
	}

	// clean directory to remove all files
	if err := postService.Filesystem.CleanDir(); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return err
	}

	// save zipped project
	if err := postService.Filesystem.SaveZipFile(c, file); err != nil {
		// it fails so we set render status to failed and reset the branch
		post.RenderStatus = models.Failure
		_, _ = postService.PostRepository.Update(post)
		_ = postService.Filesystem.Reset()

		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return fmt.Errorf("failed to save zip file: %w", err)
	}

	// commit (or perhaps only commit after rendering?)
	if err := postService.Filesystem.CreateCommit(); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return err
	}

	// Set render status pending
	post.RenderStatus = models.Pending
	if _, err := postService.PostRepository.Update(post); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return fmt.Errorf("failed to update post entity: %w", err)
	}

	go postService.RenderService.RenderPost(post, lock)

	return nil
}

func (postService *PostService) GetMainProject(postID uint) (string, *flock.Flock, error) {
	var filePath string

	// check post exists
	_, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return filePath, nil, fmt.Errorf("failed to find post with id %v: %w", postID, err)
	}

	// lock directory
	// unlock upon error or after the controller has read the zip contents
	lock, err := postService.Filesystem.LockDirectory(postID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to acquire lock for directory %v: %w", postID, err)
	}

	// select repository of the parent post
	postService.Filesystem.CheckoutDirectory(postID)

	// checkout specified branch
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return filePath, nil, fmt.Errorf("failed to find master branch: %w", err)
	}

	return postService.Filesystem.GetCurrentZipFilePath(), lock, nil
}

func (postService *PostService) GetMainFiletree(postID uint) (map[string]int64, error, error) {
	// check post exists
	_, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return nil, fmt.Errorf("failed to find post with id %v: %w", postID, err), nil
	}

	// lock directory and defer unlock
	lock, err := postService.Filesystem.LockDirectory(postID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", postID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// select repository of the parent post
	postService.Filesystem.CheckoutDirectory(postID)

	// checkout specified branch
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		return nil, fmt.Errorf("failed to find master branch: %w", err), nil
	}

	// get file tree
	fileTree, err := postService.Filesystem.GetFileTree()

	return fileTree, nil, err
}

func (postService *PostService) GetMainFileFromProject(postID uint, relFilepath string) (string, *flock.Flock, error) {
	var absFilepath string

	// validate file path is inside of repository
	if strings.Contains(relFilepath, "..") {
		return absFilepath, nil, fmt.Errorf("file is outside of repository")
	}

	// check post exists
	_, err := postService.PostRepository.GetByID(postID)

	if err != nil {
		return absFilepath, nil, fmt.Errorf("failed to find post with id %v: %w", postID, err)
	}

	// lock directory
	// unlock upon error or after controller has read file
	lock, err := postService.Filesystem.LockDirectory(postID)
	if err != nil {
		return absFilepath, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", postID, err)
	}

	// select repository of the post
	postService.Filesystem.CheckoutDirectory(postID)

	// checkout master
	if err := postService.Filesystem.CheckoutBranch("master"); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return absFilepath, nil, fmt.Errorf("failed to find master branch: %w", err)
	}

	absFilepath = filepath.Join(postService.Filesystem.GetCurrentQuartoDirPath(), relFilepath)

	// Check that file exists, if not return 404
	if exists := utils.FileExists(absFilepath); !exists {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return "", nil, fmt.Errorf("no such file exists")
	}

	return absFilepath, lock, nil
}

func (postService *PostService) Filter(page, size int) ([]uint, error) {
	// TODO construct query based off filter form
	posts, err := postService.PostRepository.QueryPaginated(page, size)
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

func (postService *PostService) GetProjectPost(postID uint) (*models.ProjectPost, error) {
	// Ensure the post itself exists
	if _, err := postService.PostRepository.GetByID(postID); err != nil {
		return nil, fmt.Errorf("failed to get post with ID %d: %w", postID, err)
	}

	// Query for a project post that has this post
	// TODO this is not super efficient... improve somehow?
	foundProjectPosts, err := postService.ProjectPostRepository.Query(fmt.Sprintf("post_id = %d", postID))
	if err != nil {
		return nil, fmt.Errorf("failed to get project post that has post ID %d: %w", postID, err)
	}

	// Ensure that only ONE project post has this post.
	// TODO this is a pretty hacky way to represent the post/project-post relation.
	numberOfFoundProjectPosts := len(foundProjectPosts)

	if numberOfFoundProjectPosts != 1 {
		return nil, fmt.Errorf("failed to get exactly 1 project post for post ID %d: found %d", postID, numberOfFoundProjectPosts)
	}

	// Guaranteed to be safe due to the above condition
	return foundProjectPosts[0], nil
}
