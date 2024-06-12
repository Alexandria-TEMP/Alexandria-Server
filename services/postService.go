package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

type PostService struct {
	PostRepository   database.ModelRepositoryInterface[*models.Post]
	MemberRepository database.ModelRepositoryInterface[*models.Member]

	PostCollaboratorService interfaces.PostCollaboratorService

	// TODO add filesystem interface
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

	postFields := form.ScientificFieldTags
	postTagContainer := tags.ScientificFieldTagContainer{
		ScientificFieldTags: postFields,
	}

	post := models.Post{
		Collaborators:               postCollaborators,
		Title:                       form.Title,
		PostType:                    form.PostType,
		ScientificFieldTagContainer: postTagContainer,
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

func (postService *PostService) Filter(page, size int, form forms.PostFilterForm) ([]uint, error) {
	// TODO construct query based off filter form
	var query string

	if form.IncludeProjectPosts {
		query = ""
	} else {
		query = "post_type != 'project'"
	}

	posts, err := postService.PostRepository.QueryPaginated(page, size, query)
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
