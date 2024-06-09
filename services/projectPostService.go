package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

type ProjectPostService struct {
	ProjectPostRepository database.ModelRepositoryInterface[*models.ProjectPost]
	MemberRepository      database.ModelRepositoryInterface[*models.Member]

	PostCollaboratorService   interfaces.PostCollaboratorService
	BranchCollaboratorService interfaces.BranchCollaboratorService
}

func (projectPostService *ProjectPostService) GetProjectPost(id uint) (*models.ProjectPost, error) {
	return projectPostService.ProjectPostRepository.GetByID(id)
}

func (projectPostService *ProjectPostService) CreateProjectPost(form *forms.ProjectPostCreationForm) (*models.ProjectPost, error) {
	// This function may only be used to create Posts of type Project.
	if form.PostCreationForm.PostType != models.Project {
		return nil, fmt.Errorf("function CreateProjectPost may only create Post of type Project. received: %s",
			form.PostCreationForm.PostType)
	}

	// Information about the creators of this Project Post
	memberIDs := form.PostCreationForm.AuthorMemberIDs
	anonymous := form.PostCreationForm.Anonymous

	postCollaborators, err := projectPostService.PostCollaboratorService.MembersToPostCollaborators(memberIDs, anonymous, models.Author)
	if err != nil {
		return nil, fmt.Errorf("could not create project post: %w", err)
	}

	// This Post instance will be embedded into the Project Post
	post := models.Post{
		Collaborators: postCollaborators,

		Title:               form.PostCreationForm.Title,
		PostType:            form.PostCreationForm.PostType,
		ScientificFieldTags: form.PostCreationForm.ScientificFieldTags,

		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}

	// A new project post starts with a single "initial proposed changes" branch. This is how project posts,
	// themselves, are initially peer reviewed. While this initial proposed changes branch is open, no other
	// branches may be opened on the Project Post.
	branchCollaborators, err := projectPostService.BranchCollaboratorService.MembersToBranchCollaborators(memberIDs, anonymous)
	if err != nil {
		return nil, fmt.Errorf("could not create project post: %w", err)
	}

	projectPost := models.ProjectPost{
		Post:               post,
		CompletionStatus:   form.CompletionStatus,
		FeedbackPreference: form.FeedbackPreference,

		OpenBranches: []*models.Branch{
			{
				// TODO make these fields optional maybe? so they dont have to be filled in
				NewPostTitle:            form.PostCreationForm.Title,
				UpdatedCompletionStatus: form.CompletionStatus,
				UpdatedScientificFields: form.PostCreationForm.ScientificFieldTags,
				Collaborators:           branchCollaborators,
				Reviews:                 []*models.BranchReview{},
				DiscussionContainer: models.DiscussionContainer{
					Discussions: []*models.Discussion{},
				},
				BranchTitle:        models.InitialPeerReviewBranchName,
				RenderStatus:       models.Pending,
				BranchReviewStatus: models.BranchOpenForReview,
			},
		},
		ClosedBranches: []*models.ClosedBranch{},

		// New project posts are always open for review
		PostReviewStatus: models.Open,
	}

	if err := projectPostService.ProjectPostRepository.Create(&projectPost); err != nil {
		return nil, fmt.Errorf("unable to create projectt post: %w", err)
	}

	return &projectPost, nil
}

func (projectPostService *ProjectPostService) UpdateProjectPost(_ *models.ProjectPost) error {
	return fmt.Errorf("TODO")
}
