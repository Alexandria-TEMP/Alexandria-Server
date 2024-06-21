package services

import (
	"fmt"
	"log"
	"slices"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

type ProjectPostService struct {
	ClosedBranchRepository                database.ModelRepositoryInterface[*models.ClosedBranch]
	PostRepository                        database.ModelRepositoryInterface[*models.Post]
	ProjectPostRepository                 database.ModelRepositoryInterface[*models.ProjectPost]
	MemberRepository                      database.ModelRepositoryInterface[*models.Member]
	ScientificFieldTagContainerRepository database.ModelRepositoryInterface[*models.ScientificFieldTagContainer]
	Filesystem                            filesystemInterfaces.Filesystem

	PostCollaboratorService   interfaces.PostCollaboratorService
	BranchCollaboratorService interfaces.BranchCollaboratorService
	BranchService             interfaces.BranchService
	TagService                interfaces.TagService
}

func (projectPostService *ProjectPostService) GetProjectPost(id uint) (*models.ProjectPost, error) {
	return projectPostService.ProjectPostRepository.GetByID(id)
}

func (projectPostService *ProjectPostService) CreateProjectPost(form *forms.ProjectPostCreationForm, member *models.Member) (*models.ProjectPost, error, error) {
	// Create the post
	post, err := projectPostService.createPostForProjectPost(form)

	if err != nil {
		return nil, fmt.Errorf("could not create post: %w", err), nil
	}

	// Information about the creators of this Project Post
	memberIDs := form.AuthorMemberIDs
	anonymous := form.Anonymous

	// check if creating member is in authors or branch is anonymous
	if !anonymous && !slices.Contains(memberIDs, member.ID) {
		return nil, fmt.Errorf("the creating member is not in the list of authors. either add the member or set the branch to anonymous"), nil
	}

	// A new project post starts with a single "initial proposed changes" branch. This is how project posts,
	// themselves, are initially peer reviewed. While this initial proposed changes branch is open, no other
	// branches may be opened on the Project Post.
	branchCollaborators, err := projectPostService.BranchCollaboratorService.MembersToBranchCollaborators(memberIDs, anonymous)
	if err != nil {
		return nil, fmt.Errorf("could not create project post: %w", err), nil
	}

	// Construct project post
	projectPost := models.ProjectPost{
		Post:                      *post,
		ProjectCompletionStatus:   form.ProjectCompletionStatus,
		ProjectFeedbackPreference: form.ProjectFeedbackPreference,

		OpenBranches: []*models.Branch{
			{
				// TODO make these fields optional maybe? so they dont have to be filled in
				Collaborators: branchCollaborators,
				Reviews:       []*models.BranchReview{},
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

	// lock directory and defer unlocking it
	lock, err := projectPostService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// Checkout directory where project post will store it's files
	directoryFilesystem := projectPostService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// Create a new git repo there
	if err := directoryFilesystem.CreateRepository(); err != nil {
		return nil, nil, err
	}

	// Create initial branch in git repo
	branch := projectPost.OpenBranches[0]
	if err := directoryFilesystem.CreateBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		return nil, nil, err
	}

	return &projectPost, nil, nil
}

func (projectPostService *ProjectPostService) createPostForProjectPost(form *forms.ProjectPostCreationForm) (*models.Post, error) {
	// Information about the creators of this Project Post
	memberIDs := form.AuthorMemberIDs
	anonymous := form.Anonymous

	postCollaborators, err := projectPostService.PostCollaboratorService.MembersToPostCollaborators(memberIDs, anonymous, models.Author)
	if err != nil {
		return nil, fmt.Errorf("failed to get post collaborators: %w", err)
	}

	// convert []uint to []*models.ScientificFieldTag
	tags, err := projectPostService.TagService.GetTagsFromIDs(form.ScientificFieldTagIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get tags from ids: %w", err)
	}

	// create and save the tag container to avoid issues with saving later (preloading stuff?)
	postTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: tags,
	}

	if err := projectPostService.ScientificFieldTagContainerRepository.Create(postTagContainer); err != nil {
		return nil, fmt.Errorf("failed to add tag container to db: %w", err)
	}

	// This Post instance will be embedded into the Project Post
	post := &models.Post{
		Collaborators:               postCollaborators,
		Title:                       form.Title,
		PostType:                    models.Project,
		ScientificFieldTagContainer: *postTagContainer,
		RenderStatus:                models.Success,
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}

	// Save post due to ScientificFieldTagContainer issues otherwise (preloading?)
	if err := projectPostService.PostRepository.Create(post); err != nil {
		return nil, fmt.Errorf("failed to add post to db: %w", err)
	}

	return post, nil
}

func (projectPostService *ProjectPostService) GetBranchesGroupedByReviewStatus(projectPostID uint) (*models.BranchesGroupedByReviewStatusDTO, error) {
	// Get the project post
	projectPost, err := projectPostService.ProjectPostRepository.GetByID(projectPostID)
	if err != nil {
		return nil, fmt.Errorf("could not find project post with ID %d: %w", projectPostID, err)
	}

	// We categorize branches in three categories:
	// 1) Open for review
	// 2) Rejected
	// 3) Approved (peer reviewed)

	openForReviewBranchIDs := make([]uint, len(projectPost.OpenBranches))
	rejectedClosedBranchIDs := []uint{}
	approvedClosedBranchIDs := []uint{}

	// Add every single open branch
	for i, branch := range projectPost.OpenBranches {
		openForReviewBranchIDs[i] = branch.ID
	}

	// Add closed branches that are rejected
	for _, branch := range projectPost.ClosedBranches {
		if branch.BranchReviewDecision == models.Rejected {
			rejectedClosedBranchIDs = append(rejectedClosedBranchIDs, branch.ID)
		}
	}

	// Add closed branches that are approved
	for _, branch := range projectPost.ClosedBranches {
		if branch.BranchReviewDecision == models.Approved {
			approvedClosedBranchIDs = append(approvedClosedBranchIDs, branch.ID)
		}
	}

	groupedBranchesByStatus := &models.BranchesGroupedByReviewStatusDTO{
		OpenBranchIDs:           openForReviewBranchIDs,
		RejectedClosedBranchIDs: rejectedClosedBranchIDs,
		ApprovedClosedBranchIDs: approvedClosedBranchIDs,
	}

	return groupedBranchesByStatus, nil
}

func (projectPostService *ProjectPostService) GetDiscussionContainersFromMergeHistory(projectPostID uint) (*models.DiscussionContainerProjectHistoryDTO, error) {
	// Get each discussion container, from every closed & merged branch. Sources of containers:
	// 1) On the underlying post
	// 2) On each closed + merged branch
	// Get the post's discussion container
	projectPost, err := projectPostService.ProjectPostRepository.GetByID(projectPostID)
	if err != nil {
		return nil, fmt.Errorf("could not get project post: %w", err)
	}

	postDiscussionContainerID := projectPost.Post.DiscussionContainerID

	// Get each closed + merged branch's discussion container
	closedApprovedBranches, err := projectPostService.ClosedBranchRepository.Query("branch_review_decision = 'approved'")
	if err != nil {
		return nil, fmt.Errorf("could not get closed approved branches of project post: %w", err)
	}

	// Transform each branch into a 'branch + discussion container DTO', which holds
	// the 'closed branch' ID, and the discussion container ID from the branch itself
	discussionContainersWithBranches := make([]models.DiscussionContainerWithBranchDTO, len(closedApprovedBranches))

	for i, closedApprovedBranch := range closedApprovedBranches {
		discussionContainersWithBranches[i] = models.DiscussionContainerWithBranchDTO{
			DiscussionContainerID: closedApprovedBranch.Branch.DiscussionContainerID,
			ClosedBranchID:        closedApprovedBranch.ID,
		}
	}

	// Create the final history DTO, holding all the discussion containers
	discussionContainerHistory := models.DiscussionContainerProjectHistoryDTO{
		CurrentDiscussionContainerID:     postDiscussionContainerID,
		MergedBranchDiscussionContainers: discussionContainersWithBranches,
	}

	return &discussionContainerHistory, nil
}
