package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

type ProjectPostService struct {
	ProjectPostRepository database.ModelRepositoryInterface[*models.ProjectPost]
	MemberRepository      database.ModelRepositoryInterface[*models.Member]
	Filesystem            filesystemInterfaces.Filesystem

	PostCollaboratorService   interfaces.PostCollaboratorService
	BranchCollaboratorService interfaces.BranchCollaboratorService
	BranchService             interfaces.BranchService
}

func (projectPostService *ProjectPostService) GetProjectPost(id uint) (*models.ProjectPost, error) {
	return projectPostService.ProjectPostRepository.GetByID(id)
}

func (projectPostService *ProjectPostService) CreateProjectPost(form *forms.ProjectPostCreationForm) (*models.ProjectPost, error, error) {
	// This function may only be used to create Posts of type Project.
	if form.PostCreationForm.PostType != models.Project {
		return nil, fmt.Errorf("function CreateProjectPost may only create Post of type Project. received: %s",
			form.PostCreationForm.PostType), nil
	}

	// Information about the creators of this Project Post
	memberIDs := form.PostCreationForm.AuthorMemberIDs
	anonymous := form.PostCreationForm.Anonymous

	postCollaborators, err := projectPostService.PostCollaboratorService.MembersToPostCollaborators(memberIDs, anonymous, models.Author)
	if err != nil {
		return nil, fmt.Errorf("could not create project post: %w", err), nil
	}

	// This Post instance will be embedded into the Project Post
	post := models.Post{
		Collaborators: postCollaborators,

		Title:            form.PostCreationForm.Title,
		PostType:         form.PostCreationForm.PostType,
		ScientificFields: form.PostCreationForm.ScientificFields,
		RenderStatus:     models.Success,

		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}

	// A new project post starts with a single "initial proposed changes" branch. This is how project posts,
	// themselves, are initially peer reviewed. While this initial proposed changes branch is open, no other
	// branches may be opened on the Project Post.
	branchCollaborators, err := projectPostService.BranchCollaboratorService.MembersToBranchCollaborators(memberIDs, anonymous)
	if err != nil {
		return nil, fmt.Errorf("could not create project post: %w", err), nil
	}

	projectPost := models.ProjectPost{
		Post:                      post,
		ProjectCompletionStatus:   form.ProjectCompletionStatus,
		ProjectFeedbackPreference: form.ProjectFeedbackPreference,

		OpenBranches: []*models.Branch{
			{
				// TODO make these fields optional maybe? so they dont have to be filled in
				UpdatedPostTitle:        &form.PostCreationForm.Title,
				UpdatedCompletionStatus: &form.ProjectCompletionStatus,
				UpdatedScientificFields: form.PostCreationForm.ScientificFields,
				Collaborators:           branchCollaborators,
				Reviews:                 []*models.BranchReview{},
				DiscussionContainer: models.DiscussionContainer{
					Discussions: []*models.Discussion{},
				},
				BranchTitle:               models.InitialPeerReviewBranchName,
				RenderStatus:              models.Success,
				BranchOverallReviewStatus: models.BranchOpenForReview,
			},
		},
		ClosedBranches: []*models.ClosedBranch{},

		// New project posts are always open for branchreview
		PostReviewStatus: models.Open,
	}

	// Add the project post to db
	if err := projectPostService.ProjectPostRepository.Create(&projectPost); err != nil {
		return nil, nil, fmt.Errorf("unable to create project post: %w", err)
	}

	// Checkout directory where project post will store it's files
	projectPostService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// Create a new git repo there
	if err := projectPostService.Filesystem.CreateRepository(); err != nil {
		return nil, nil, err
	}

	// Create initial branch in git repo
	branch := projectPost.OpenBranches[0]
	if err := projectPostService.Filesystem.CreateBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		return nil, nil, err
	}

	return &projectPost, nil, nil
}

func (projectPostService *ProjectPostService) UpdateProjectPost(_ *models.ProjectPost) error {
	return fmt.Errorf("TODO")
}

func (projectPostService *ProjectPostService) Filter(page, size int, _ forms.FilterForm) ([]uint, error) {
	// TODO construct query based off filter form
	posts, err := projectPostService.ProjectPostRepository.QueryPaginated(page, size)
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
