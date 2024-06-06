package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

type ProjectPostService struct {
	ProjectPostRepository database.ModelRepositoryInterface[*models.ProjectPost]
	MemberRepository      database.ModelRepositoryInterface[*models.Member]

	PostCollaboratorService interfaces.PostCollaboratorService
}

func (projectPostService *ProjectPostService) GetProjectPost(id uint) (*models.ProjectPost, error) {
	return projectPostService.ProjectPostRepository.GetByID(id)
}

func (projectPostService *ProjectPostService) CreateProjectPost(form *forms.ProjectPostCreationForm) (*models.ProjectPost, error) {
	// This function may only be used to create Posts of type Project.
	if form.PostCreationForm.PostType != tags.Project {
		return nil, fmt.Errorf("function CreateProjectPost may only create Post of type Project. received: %s",
			form.PostCreationForm.PostType)
	}

	memberIDs := form.PostCreationForm.AuthorMemberIDs
	anonymous := form.PostCreationForm.Anonymous

	postCollaborators, err := projectPostService.PostCollaboratorService.MembersToPostCollaborators(memberIDs, anonymous, models.Author)
	if err != nil {
		return nil, fmt.Errorf("could not create project post: %w", err)
	}

	post := models.Post{
		Collaborators: postCollaborators,

		Title:               form.PostCreationForm.Title,
		PostType:            form.PostCreationForm.PostType,
		ScientificFieldTags: form.PostCreationForm.ScientificFieldTags,

		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}

	projectPost := models.ProjectPost{
		Post:               post,
		CompletionStatus:   form.CompletionStatus,
		FeedbackPreference: form.FeedbackPreference,

		// New project posts have no open or closed branches
		// TODO change this to have the initial peer review branch
		OpenBranches:   []*models.Branch{},
		ClosedBranches: []*models.ClosedBranch{},

		// New project posts are always open for review
		PostReviewStatusTag: tags.Open,
	}

	if err := projectPostService.ProjectPostRepository.Create(&projectPost); err != nil {
		return nil, fmt.Errorf("unable to create projectt post: %w", err)
	}

	return &projectPost, nil
}

func (projectPostService *ProjectPostService) UpdateProjectPost(_ *models.ProjectPost) error {
	return fmt.Errorf("TODO")
}
