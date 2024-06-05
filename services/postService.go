package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type PostService struct {
	PostRepository   database.ModelRepositoryInterface[*models.Post]
	MemberRepository database.ModelRepositoryInterface[*models.Member]
}

func (postService *PostService) GetPost(id uint) (*models.Post, error) {
	return postService.PostRepository.GetByID(id)
}

func (postService *PostService) CreatePost(form *forms.PostCreationForm) (*models.Post, error) {
	// Posts created via this function may not be project posts
	// (those must use ProjectPostCreationForms)
	if form.PostType == tags.Project {
		return nil, fmt.Errorf("creating post of type ProjectPost using CreatePost is forbidden")
	}

	// If the Post is not anonymous, convert the author member IDs into a list of Post Collaborators
	postCollaborators := make([]*models.PostCollaborator, 0)

	// TODO extract into separate function
	if !form.Anonymous {
		authorMemberIDs := form.AuthorMemberIDs
		postCollaborators = make([]*models.PostCollaborator, len(authorMemberIDs))

		for i, memberID := range authorMemberIDs {
			// Fetch the member from the database
			member, err := postService.MemberRepository.GetByID(memberID)
			if err != nil {
				return nil, fmt.Errorf("could not create post collaborators: %w", err)
			}

			newPostCollaborator := models.PostCollaborator{
				Member: *member,
				// Post Collaborators on a brand new Post are always automatically set to Authors
				CollaborationType: models.Author,
			}

			postCollaborators[i] = &newPostCollaborator
		}
	}

	post := models.Post{
		Collaborators:       postCollaborators,
		Title:               form.Title, // TODO sanitize title?
		PostType:            form.PostType,
		ScientificFieldTags: form.ScientificFieldTags,
		DiscussionContainer: models.DiscussionContainer{
			// The discussion list is initially empty
			Discussions: []*models.Discussion{},
		},
	}

	if err := postService.PostRepository.Create(&post); err != nil {
		return nil, err
	}

	return &post, nil
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
