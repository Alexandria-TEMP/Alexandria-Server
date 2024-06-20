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

const approvalsToMerge = 2 // 0 indexed

type BranchService struct {
	BranchRepository                   database.ModelRepositoryInterface[*models.Branch]
	ClosedBranchRepository             database.ModelRepositoryInterface[*models.ClosedBranch]
	PostRepository                     database.ModelRepositoryInterface[*models.Post]
	ProjectPostRepository              database.ModelRepositoryInterface[*models.ProjectPost]
	ReviewRepository                   database.ModelRepositoryInterface[*models.BranchReview]
	DiscussionContainerRepository      database.ModelRepositoryInterface[*models.DiscussionContainer]
	DiscussionRepository               database.ModelRepositoryInterface[*models.Discussion]
	MemberRepository                   database.ModelRepositoryInterface[*models.Member]
	ScientificFieldTagRepository       database.ModelRepositoryInterface[*models.ScientificFieldTag]
	Filesystem                         filesystemInterfaces.Filesystem
	ScientificFieldTagContainerService interfaces.ScientificFieldTagContainerService

	RenderService             interfaces.RenderService
	BranchCollaboratorService interfaces.BranchCollaboratorService
	PostCollaboratorService   interfaces.PostCollaboratorService
	TagService                interfaces.TagService
}

func (branchService *BranchService) GetBranch(branchID uint) (*models.Branch, error) {
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	return branch, nil
}

func (branchService *BranchService) CreateBranch(branchCreationForm *forms.BranchCreationForm, member *models.Member) (*models.Branch, error, error) {
	// verify parent project post exists
	projectPost, err := branchService.ProjectPostRepository.GetByID(branchCreationForm.ProjectPostID)

	if err != nil {
		return nil, fmt.Errorf("failed to find project post with id %v: %w", branchCreationForm.ProjectPostID, err), nil
	}

	// check whether project post is still open. if so reject this branch creation
	if projectPost.PostReviewStatus == models.Open {
		return nil, fmt.Errorf("this project post is still open for review"), nil
	}

	// check if creating member is in collaborators or branch is anonymous
	if !branchCreationForm.Anonymous && !slices.Contains(branchCreationForm.CollaboratingMemberIDs, member.ID) {
		return nil, fmt.Errorf("the creating member is not in the list of collaborators. either add the member or set the branch to anonymous"), nil
	}

	// create and save discussion new container
	// we shouldn't have to do this extra, it should be implicit but it isnt...
	discussionContainer := models.DiscussionContainer{}
	if err := branchService.DiscussionContainerRepository.Create(&discussionContainer); err != nil {
		return nil, fmt.Errorf("failed to add discussion container to db: %w", err), nil
	}

	// get all collaborators from ids
	collaborators, err := branchService.BranchCollaboratorService.MembersToBranchCollaborators(branchCreationForm.CollaboratingMemberIDs, branchCreationForm.Anonymous)
	if err != nil {
		return nil, fmt.Errorf("failed to convert member ids to branch collaborators: %w", err), nil
	}

	// convert []uint to []*models.ScientificFieldTag
	tags, err := branchService.TagService.GetTagsFromIDs(branchCreationForm.UpdatedScientificFieldIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get tags from ids: %w", err), nil
	}

	// make new branch
	branch := &models.Branch{
		UpdatedPostTitle:                   branchCreationForm.UpdatedPostTitle,
		UpdatedCompletionStatus:            branchCreationForm.UpdatedCompletionStatus,
		UpdatedFeedbackPreferences:         branchCreationForm.UpdatedFeedbackPreferences,
		UpdatedScientificFieldTagContainer: &models.ScientificFieldTagContainer{ScientificFieldTags: tags},
		Collaborators:                      collaborators,
		DiscussionContainer:                discussionContainer,
		ProjectPostID:                      &branchCreationForm.ProjectPostID,
		BranchTitle:                        branchCreationForm.BranchTitle,
		RenderStatus:                       models.Success,
		BranchOverallReviewStatus:          models.BranchOpenForReview,
	}

	// save branch entity to open branches
	projectPost.OpenBranches = append(projectPost.OpenBranches, branch)

	if _, err := branchService.ProjectPostRepository.Update(projectPost); err != nil {
		return nil, nil, fmt.Errorf("failed to update project post with new branch: %w", err)
	}

	// lock directory and defer unlocking it
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// set vfs to repository according to the Post of the ProjectPost of the Branch entity
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// create new branch in git repo with branch ID as its name
	if err := branchService.Filesystem.CreateBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		return nil, nil, fmt.Errorf("failed create branch: %w", err)
	}

	return branch, nil, nil
}

func (branchService *BranchService) DeleteBranch(branchID uint) error {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return fmt.Errorf("failed to get the project post of branch %v: %w", branch.ID, err)
	}

	// lock directory and defer unlocking it
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// checkout repository
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// delete branch
	if err := branchService.Filesystem.DeleteBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return fmt.Errorf("failed to delete branch from vfs with id %v: %w", branchID, err)
	}

	// delete entity
	if err := branchService.BranchRepository.Delete(branchID); err != nil {
		return fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	return nil
}

func (branchService *BranchService) GetAllBranchReviewStatuses(branchID uint) ([]models.BranchReviewDecision, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get all decisions
	decisions := []models.BranchReviewDecision{}
	for _, branchreview := range branch.Reviews {
		decisions = append(decisions, branchreview.BranchReviewDecision)
	}

	return decisions, nil
}

func (branchService *BranchService) GetReview(reviewID uint) (*models.BranchReview, error) {
	// get branch
	branchreview, err := branchService.ReviewRepository.GetByID(reviewID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v: %w", reviewID, err)
	}

	return branchreview, nil
}

func (branchService *BranchService) CreateReview(form forms.ReviewCreationForm, member *models.Member) (*models.BranchReview, error) {
	if canReview, err401, err404 := branchService.MemberCanReview(form.BranchID, member); !canReview {
		return nil, fmt.Errorf("this member cannot review this branch: %w", err401)
	} else if err404 != nil {
		return nil, fmt.Errorf("failed to check whether this member can review the branch: %w", err404)
	}

	// get branch
	branch, err := branchService.BranchRepository.GetByID(form.BranchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v: %w", form.BranchID, err)
	}

	// ensure the branch isn't already closed
	if branch.BranchOverallReviewStatus != models.BranchOpenForReview {
		return nil, fmt.Errorf("branch is already reviewed with status '%v'", branch.BranchOverallReviewStatus)
	}

	// make new branchreview
	branchreview := &models.BranchReview{
		BranchID:             form.BranchID,
		Member:               *member,
		BranchReviewDecision: form.BranchReviewDecision,
		Feedback:             form.Feedback,
	}

	if err := branchService.ReviewRepository.Create(branchreview); err != nil {
		return nil, fmt.Errorf("failed to add branch review to db: %w", err)
	}

	// update branch with new branchreview and update branchreview status accordingly
	branch.Reviews = append(branch.Reviews, branchreview)
	branch.BranchOverallReviewStatus = branchService.updateReviewStatus(branch.Reviews)

	// if approved or rejected we close the branch
	if branch.BranchOverallReviewStatus == models.BranchPeerReviewed || branch.BranchOverallReviewStatus == models.BranchRejected {
		if err := branchService.closeBranch(branch); err != nil {
			return nil, fmt.Errorf("failed to close branch: %w", err)
		}

		return branchreview, nil
	}

	// save changes to branch
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return nil, fmt.Errorf("failed to save branch branchreview: %w", err)
	}

	return branchreview, nil
}

func (branchService *BranchService) closeBranch(branch *models.Branch) error {
	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return fmt.Errorf("failed to get the project post of branch %v: %w", branch.ID, err)
	}

	// close branch
	closedBranch := &models.ClosedBranch{
		ProjectPostID: projectPost.ID,
	}

	// merge into master if approved
	if branch.BranchOverallReviewStatus == models.BranchPeerReviewed {
		closedBranch.BranchReviewDecision = models.Approved
		projectPost.PostReviewStatus = models.Reviewed

		if err := branchService.merge(branch, closedBranch, projectPost); err != nil {
			return err
		}
	} else {
		closedBranch.BranchReviewDecision = models.Rejected

		// If the branch was rejected, and the project post was "open for review", this
		// means the project post itself has been marked "revision needed".
		// If the branch was rejected, and the project post was already peer reviewed,
		// it shall remain peer reviewed.
		if projectPost.PostReviewStatus == models.Open {
			projectPost.PostReviewStatus = models.RevisionNeeded
		}
	}

	// remove project post id so that it is no longer in open branches
	branch.ProjectPostID = nil
	closedBranch.Branch = *branch

	// add branch to closed branches
	projectPost.ClosedBranches = append(projectPost.ClosedBranches, closedBranch)

	// remove branch from open branches
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	newOpenBranches := []*models.Branch{}

	for _, b := range projectPost.OpenBranches {
		if b.ID != branch.ID {
			newOpenBranches = append(newOpenBranches, b)
		}
	}

	projectPost.OpenBranches = newOpenBranches

	// save changes to project post and branch
	if _, err := branchService.ProjectPostRepository.Update(projectPost); err != nil {
		return fmt.Errorf("failed to save project post and branch: %w", err)
	}

	return nil
}

func (branchService *BranchService) merge(branch *models.Branch, closedBranch *models.ClosedBranch, projectPost *models.ProjectPost) error {
	closedBranch.BranchReviewDecision = models.Approved

	// lock directory and defer unlocking it
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// checkout repo and then merge
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	if err := branchService.Filesystem.Merge(fmt.Sprintf("%v", branch.ID), "master"); err != nil {
		return err
	}

	// find the last branch merged into this project post
	mergedBranches, err := branchService.ClosedBranchRepository.Query(&models.ClosedBranch{
		ProjectPostID:        projectPost.ID,
		BranchReviewDecision: models.Approved,
	})

	if err != nil {
		return fmt.Errorf("failed to find merged branches in ClosedBranchRepository: %w", err)
	}

	if len(mergedBranches) >= 1 {
		closedBranch.SupercededBranch = &mergedBranches[0].Branch
	}

	// merge metadata updates to the project post
	if branch.UpdatedPostTitle != nil {
		projectPost.Post.Title = *branch.UpdatedPostTitle
	}

	if branch.UpdatedCompletionStatus != nil {
		projectPost.ProjectCompletionStatus = *branch.UpdatedCompletionStatus
	}

	if branch.UpdatedScientificFieldTagContainer != nil {
		container, err := branchService.ScientificFieldTagContainerService.GetScientificFieldTagContainer(*branch.UpdatedScientificFieldTagContainerID)
		if err != nil {
			return fmt.Errorf("could not get scientific field tag container during merge: %w", err)
		}

		projectPost.Post.ScientificFieldTagContainer = *container
	}

	if branch.UpdatedFeedbackPreferences != nil {
		projectPost.ProjectFeedbackPreference = *branch.UpdatedFeedbackPreferences
	}

	// update project post contributors
	if err := branchService.PostCollaboratorService.MergeContributors(projectPost, branch.Collaborators); err != nil {
		return err
	}

	// update project post reviewers
	if err := branchService.PostCollaboratorService.MergeReviewers(projectPost, branch.Reviews); err != nil {
		return err
	}

	// update the post itself
	if _, err := branchService.PostRepository.Update(&projectPost.Post); err != nil {
		return fmt.Errorf("could not update post: %w", err)
	}

	return nil
}

// updateReviewStatus finds the current branchreview status
// If there are 3 approvals, approve the branch
// If there are any rejections, reject the branch
// Otherwise leave open for branchreview
func (branchService *BranchService) updateReviewStatus(reviews []*models.BranchReview) models.BranchOverallReviewStatus {
	for i, r := range reviews {
		if r.BranchReviewDecision == models.Rejected {
			return models.BranchRejected
		}

		if i == approvalsToMerge {
			return models.BranchPeerReviewed
		}
	}

	return models.BranchOpenForReview
}

func (branchService *BranchService) MemberCanReview(branchID uint, member *models.Member) (bool, error, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// check if branch is open for review
	if branch.BranchOverallReviewStatus != models.BranchOpenForReview {
		return false, fmt.Errorf("this branch has been closed and cannot be reviewed"), nil
	}

	// check if member is contributor to branch
	branchCollaboratorMemberIDs := make([]uint, 0)
	for _, branchCollaborator := range branch.Collaborators {
		branchCollaboratorMemberIDs = append(branchCollaboratorMemberIDs, branchCollaborator.MemberID)
	}

	if slices.Contains(branchCollaboratorMemberIDs, member.ID) {
		return false, fmt.Errorf("this member has contributed to the branch, so they cannot review it"), nil
	}

	// check if member has already reviewed branch
	reviewMemberIDs := make([]uint, 0)
	for _, review := range branch.Reviews {
		reviewMemberIDs = append(reviewMemberIDs, review.MemberID)
	}

	if slices.Contains(reviewMemberIDs, member.ID) {
		return false, fmt.Errorf("this member has already reviewed this branch, so they cannot review it"), nil
	}

	// get the tag ids of the branch.
	branchTagIDs, err := branchService.getBranchTagIDs(branch)

	if err != nil {
		return false, nil, fmt.Errorf("failed to find branch tags: %w", err)
	}

	// if not tags then always true
	if len(branchTagIDs) == 0 {
		return true, nil, nil
	}

	// get scientific field tags of member
	memberTagsContainer, err := branchService.ScientificFieldTagContainerService.GetScientificFieldTagContainer(member.ScientificFieldTagContainerID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to find scientific field tag container of member: %w", err)
	}

	memberTagIDs := models.ScientificFieldTagIntoIDs(memberTagsContainer.ScientificFieldTags)

	// iterate over each member tag to see if there is an intersection with branch tags
	for _, memberTagID := range memberTagIDs {
		if slices.Contains(branchTagIDs, memberTagID) {
			return true, nil, nil
		}
	}

	// if no intersection is found we return false
	return false, fmt.Errorf("this member is not qualified to review this branch"), nil
}

func (branchService *BranchService) getBranchTagIDs(branch *models.Branch) ([]uint, error) {
	// if the branch has updated the tags, we use those, if the branch has not then we use the posts tags.
	var branchTags []*models.ScientificFieldTag

	if branch.UpdatedScientificFieldTagContainerID != nil {
		container, err := branchService.ScientificFieldTagContainerService.GetScientificFieldTagContainer(*branch.UpdatedScientificFieldTagContainerID)
		if err != nil {
			return nil, fmt.Errorf("failed to find scientific field tag container of branch: %w", err)
		}

		branchTags = container.ScientificFieldTags
	} else {
		projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)
		if err != nil {
			return nil, fmt.Errorf("failed to find project post of branch: %w", err)
		}

		container, err := branchService.ScientificFieldTagContainerService.GetScientificFieldTagContainer(projectPost.Post.ScientificFieldTagContainerID)
		if err != nil {
			return nil, fmt.Errorf("failed to find scientific field tag container of post: %w", err)
		}

		branchTags = container.ScientificFieldTags
	}

	// recursively add all of the parents of the tags to the list of branchTags
	// because if a member is qualified in a field that subsumes a branchTag, they should be able to review it.
	for i := 0; i < len(branchTags); i++ {
		branchTag := branchTags[i]

		if branchTag.ParentID == nil {
			continue
		}

		parentTag, err := branchService.ScientificFieldTagRepository.GetByID(*branchTag.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to find parent scientific field tag: %w", err)
		} else if !slices.Contains(branchTags, parentTag) {
			branchTags = append(branchTags, parentTag)
		}
	}

	// map all the branch tags to their IDs
	return models.ScientificFieldTagIntoIDs(branchTags), nil
}

func (branchService *BranchService) GetProject(branchID uint) (string, *flock.Flock, error) {
	var filePath string

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return "", nil, err
	}

	// lock directory.
	// we unlock in the controller once the project file has been read or if we error.
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return filePath, nil, fmt.Errorf("failed to find this git branch, with name %v: %w", branchID, err)
	}

	return branchService.Filesystem.GetCurrentZipFilePath(), lock, nil
}

func (branchService *BranchService) UploadProject(c *gin.Context, file *multipart.FileHeader, branchID uint) error {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return err
	}

	// lock directory
	// if there is an error we will unlock, otherwise we unlock at the end of the render pipeline
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return err
	}

	// clean directory to remove all files
	if err := branchService.Filesystem.CleanDir(); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return err
	}

	// save zipped project
	if err := branchService.Filesystem.SaveZipFile(c, file); err != nil {
		// it fails so we set render status to failed and reset the branch
		branch.RenderStatus = models.Failure
		_, _ = branchService.BranchRepository.Update(branch)
		_ = branchService.Filesystem.Reset()

		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return fmt.Errorf("failed to save zip file: %w", err)
	}

	// commit
	if err := branchService.Filesystem.CreateCommit(); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return err
	}

	// Set render status pending
	branch.RenderStatus = models.Pending
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return fmt.Errorf("failed to update branch entity: %w", err)
	}

	go branchService.RenderService.RenderBranch(branch, lock)

	return nil
}

func (branchService *BranchService) GetFiletree(branchID uint) (map[string]int64, error, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err), nil
	}

	// get project post. if branch is clsoed we need to get the project post id via the closed branch
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return nil, err, nil
	}

	// lock directory and defer unlocking it
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return nil, fmt.Errorf("failed to find this git branch, with name %v", branchID), nil
	}

	// get file tree
	fileTree, err := branchService.Filesystem.GetFileTree()

	return fileTree, nil, err
}

func (branchService *BranchService) GetBranchProjectPost(branch *models.Branch) (*models.ProjectPost, error) {
	var projectPostID uint

	if branch.ProjectPostID == nil {
		closedBranches, err := branchService.ClosedBranchRepository.Query(&models.ClosedBranch{BranchID: branch.ID})
		if err != nil || len(closedBranches) == 0 {
			return nil, fmt.Errorf("failed to find the closed branch for branch with id %v: %w", branch.ID, err)
		}

		projectPostID = closedBranches[0].ProjectPostID
	} else {
		projectPostID = *branch.ProjectPostID
	}

	projectPost, err := branchService.ProjectPostRepository.GetByID(projectPostID)

	if err != nil {
		return nil, fmt.Errorf("failed to get project post with id %v: %w", projectPostID, err)
	}

	// set discussion container (isn't preloaded properly)
	discussionContainer, err := branchService.DiscussionContainerRepository.GetByID(projectPost.Post.DiscussionContainerID)

	if err != nil {
		return nil, fmt.Errorf("failed to get discussion container: %w", err)
	}

	projectPost.Post.DiscussionContainer = *discussionContainer

	// set closed branches (isn't preloaded properly)
	projectPost.ClosedBranches, err = branchService.ClosedBranchRepository.Query(&models.ClosedBranch{ProjectPostID: projectPost.ID})

	if err != nil {
		return nil, fmt.Errorf("failed to get closed branches for project post: %w", err)
	}

	return projectPost, nil
}

func (branchService *BranchService) GetFileFromProject(branchID uint, relFilepath string) (string, *flock.Flock, error) {
	var absFilepath string

	// validate file path is inside of repository
	if strings.Contains(relFilepath, "..") {
		return absFilepath, nil, fmt.Errorf("file is outside of repository")
	}

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return absFilepath, nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return "", nil, err
	}

	// lock directory
	// we unlock in the controller after the file has been read from the reposioptory or if there is an error
	lock, err := branchService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return absFilepath, nil, fmt.Errorf("failed to find this git branch, with name %v: %w", branchID, err)
	}

	absFilepath = filepath.Join(branchService.Filesystem.GetCurrentQuartoDirPath(), relFilepath)

	// Check that file exists, if not return 404
	if exists := utils.FileExists(absFilepath); !exists {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return "", nil, fmt.Errorf("no such file exists")
	}

	return absFilepath, lock, nil
}

func (branchService *BranchService) GetClosedBranch(closedBranchID uint) (*models.ClosedBranch, error) {
	closedBranch, err := branchService.ClosedBranchRepository.GetByID(closedBranchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find closed branch with id %v: %w", closedBranchID, err)
	}

	return closedBranch, nil
}
