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

const approvalsToMerge = 2 // 0 indexed

type BranchService struct {
	BranchRepository              database.ModelRepositoryInterface[*models.Branch]
	ClosedBranchRepository        database.ModelRepositoryInterface[*models.ClosedBranch]
	PostRepository                database.ModelRepositoryInterface[*models.Post]
	ProjectPostRepository         database.ModelRepositoryInterface[*models.ProjectPost]
	ReviewRepository              database.ModelRepositoryInterface[*models.BranchReview]
	DiscussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
	DiscussionRepository          database.ModelRepositoryInterface[*models.Discussion]
	MemberRepository              database.ModelRepositoryInterface[*models.Member]
	Filesystem                    filesystemInterfaces.Filesystem

	RenderService             interfaces.RenderService
	BranchCollaboratorService interfaces.BranchCollaboratorService
	PostCollaboratorService   interfaces.PostCollaboratorService
	TagService                interfaces.TagService
}

func (branchService *BranchService) GetBranch(branchID uint) (models.Branch, error) {
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return *branch, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	return *branch, nil
}

func (branchService *BranchService) CreateBranch(branchCreationForm *forms.BranchCreationForm) (models.Branch, error, error) {
	var branch models.Branch

	// verify parent project post exists
	projectPost, err := branchService.ProjectPostRepository.GetByID(branchCreationForm.ProjectPostID)

	if err != nil {
		return branch, fmt.Errorf("failed to find project post with id %v: %w", branchCreationForm.ProjectPostID, err), nil
	}

	// create and save discussion new container
	// we shouldn't have to do this extra, it should be implicit but it isnt...
	discussionContainer := models.DiscussionContainer{}
	if err := branchService.DiscussionContainerRepository.Create(&discussionContainer); err != nil {
		return branch, fmt.Errorf("failed to add discussion container to db: %w", err), nil
	}

	// get all collaborators from ids
	collaborators, err := branchService.BranchCollaboratorService.MembersToBranchCollaborators(branchCreationForm.CollaboratingMemberIDs, branchCreationForm.Anonymous)
	if err != nil {
		return branch, fmt.Errorf("failed to convert member ids to branch collaborators: %w", err), nil
	}

	// convert []uint to []*models.ScientificFieldTag
	tags, err := branchService.TagService.GetTagsFromIDs(branchCreationForm.UpdatedScientificFieldIDs)

	if err != nil {
		return branch, fmt.Errorf("failed to get tags from ids: %w", err), nil
	}

	// make new branch
	branch = models.Branch{
		UpdatedPostTitle:                   branchCreationForm.UpdatedPostTitle,
		UpdatedCompletionStatus:            branchCreationForm.UpdatedCompletionStatus,
		UpdatedScientificFieldTagContainer: &models.ScientificFieldTagContainer{ScientificFieldTags: tags},
		Collaborators:                      collaborators,
		DiscussionContainer:                discussionContainer,
		ProjectPostID:                      &branchCreationForm.ProjectPostID,
		BranchTitle:                        branchCreationForm.BranchTitle,
		RenderStatus:                       models.Success,
		BranchOverallReviewStatus:          models.BranchOpenForReview,
	}

	// save branch entity to open branches
	projectPost.OpenBranches = append(projectPost.OpenBranches, &branch)

	if _, err := branchService.ProjectPostRepository.Update(projectPost); err != nil {
		return branch, nil, fmt.Errorf("failed to update project post with new branch: %w", err)
	}

	// set vfs to repository according to the Post of the ProjectPost of the Branch entity
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// create new branch in git repo with branch ID as its name
	if err := branchService.Filesystem.CreateBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		return branch, nil, fmt.Errorf("failed create branch: %w", err)
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

func (branchService *BranchService) GetReviewStatus(branchID uint) ([]models.BranchReviewDecision, error) {
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

func (branchService *BranchService) GetReview(reviewID uint) (models.BranchReview, error) {
	// get branch
	branchreview, err := branchService.ReviewRepository.GetByID(reviewID)

	if err != nil {
		return models.BranchReview{}, fmt.Errorf("failed to find branch with id %v: %w", reviewID, err)
	}

	return *branchreview, nil
}

func (branchService *BranchService) CreateReview(form forms.ReviewCreationForm) (models.BranchReview, error) {
	var branchreview models.BranchReview

	// get branch
	branch, err := branchService.BranchRepository.GetByID(form.BranchID)

	if err != nil {
		return branchreview, fmt.Errorf("failed to find branch with id %v: %w", form.BranchID, err)
	}

	// get member
	member, err := branchService.MemberRepository.GetByID(form.ReviewingMemberID)

	if err != nil {
		return branchreview, fmt.Errorf("failed to find member with id %v: %w", form.ReviewingMemberID, err)
	}

	// make new branchreview
	branchreview = models.BranchReview{
		BranchID:             form.BranchID,
		Member:               *member,
		BranchReviewDecision: form.BranchReviewDecision,
		Feedback:             form.Feedback,
	}

	if err := branchService.ReviewRepository.Create(&branchreview); err != nil {
		return branchreview, fmt.Errorf("failed to add branch review to db: %w", err)
	}

	// update branch with new branchreview and update branchreview status accordingly
	branch.Reviews = append(branch.Reviews, &branchreview)
	branch.BranchOverallReviewStatus = branchService.updateReviewStatus(branch.Reviews)

	// if approved or rejected we close the branch
	if branch.BranchOverallReviewStatus == models.BranchPeerReviewed || branch.BranchOverallReviewStatus == models.BranchRejected {
		if err := branchService.closeBranch(branch); err != nil {
			return branchreview, fmt.Errorf("failed to close branch: %w", err)
		}

		return branchreview, nil
	}

	// save changes to branch
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return branchreview, fmt.Errorf("failed to save branch branchreview: %w", err)
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
		ProjectPostID:        projectPost.ID,
		BranchReviewDecision: models.Rejected,
	}

	// merge into master if approved
	if branch.BranchOverallReviewStatus == models.BranchPeerReviewed {
		if err := branchService.merge(branch, closedBranch, projectPost); err != nil {
			return err
		}
	}

	// remove project post id so that it is no longer in open branches
	branch.ProjectPostID = nil
	closedBranch.Branch = *branch

	// add branch to closed branches
	projectPost.ClosedBranches = append(projectPost.ClosedBranches, closedBranch)

	// remove branch from open branches
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return fmt.Errorf("failed to update branch")
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
		return fmt.Errorf("failed to save project post and branch")
	}

	return nil
}

func (branchService *BranchService) merge(branch *models.Branch, closedBranch *models.ClosedBranch, projectPost *models.ProjectPost) error {
	closedBranch.BranchReviewDecision = models.Approved

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
		return fmt.Errorf("failed to find merged branches in ClosedBranchRepository")
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
		projectPost.Post.ScientificFieldTagContainer = *branch.UpdatedScientificFieldTagContainer
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

	// save changes to post (isn't being saved properly at the end for some reason)
	if _, err := branchService.PostRepository.Update(&projectPost.Post); err != nil {
		return fmt.Errorf("failed to update post metadata")
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

func (branchService *BranchService) MemberCanReview(_, _ uint) (bool, error) {
	return false, nil
}

func (branchService *BranchService) GetProject(branchID uint) (string, error) {
	var filePath string

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return "", err
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return filePath, fmt.Errorf("failed to find this git branch, with name %v", branchID)
	}

	return branchService.Filesystem.GetCurrentZipFilePath(), nil
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

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return err
	}

	// clean directory to remove all files
	if err := branchService.Filesystem.CleanDir(); err != nil {
		return err
	}

	// save zipped project
	if err := branchService.Filesystem.SaveZipFile(c, file); err != nil {
		// it fails so we set render status to failed and reset the branch
		branch.RenderStatus = models.Failure
		_, _ = branchService.BranchRepository.Update(branch)
		_ = branchService.Filesystem.Reset()

		return fmt.Errorf("failed to save zip file")
	}

	// commit
	if err := branchService.Filesystem.CreateCommit(); err != nil {
		return err
	}

	// Set render status pending
	branch.RenderStatus = models.Pending
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return fmt.Errorf("failed to update branch entity")
	}

	go branchService.RenderService.RenderBranch(branch)

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
			return nil, fmt.Errorf("failed to find the closed branch for branch with id %v", branch.ID)
		}

		projectPostID = closedBranches[0].ProjectPostID
	} else {
		projectPostID = *branch.ProjectPostID
	}

	projectPost, err := branchService.ProjectPostRepository.GetByID(projectPostID)

	if err != nil {
		return nil, fmt.Errorf("failed to get project post with id %v", projectPostID)
	}

	// set discussion container (isn't preloaded properly)
	discussionContainer, err := branchService.DiscussionContainerRepository.GetByID(projectPost.Post.DiscussionContainerID)

	if err != nil {
		return nil, fmt.Errorf("failed to get discussion container")
	}

	projectPost.Post.DiscussionContainer = *discussionContainer

	// set closed branches (isn't preloaded properly)
	projectPost.ClosedBranches, err = branchService.ClosedBranchRepository.Query(&models.ClosedBranch{ProjectPostID: projectPost.ID})

	if err != nil {
		return nil, fmt.Errorf("failed to get closed branches for project post")
	}

	return projectPost, nil
}

func (branchService *BranchService) GetFileFromProject(branchID uint, relFilepath string) (string, error) {
	var absFilepath string

	// validate file path is inside of repository
	if strings.Contains(relFilepath, "..") {
		return absFilepath, fmt.Errorf("file is outside of repository")
	}

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return absFilepath, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := branchService.GetBranchProjectPost(branch)

	if err != nil {
		return "", err
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return absFilepath, fmt.Errorf("failed to find this git branch, with name %v", branchID)
	}

	absFilepath = filepath.Join(branchService.Filesystem.GetCurrentQuartoDirPath(), relFilepath)

	// Check that file exists, if not return 404
	if exists := utils.FileExists(absFilepath); !exists {
		return "", fmt.Errorf("no such file exists")
	}

	return absFilepath, nil
}

func (branchService *BranchService) GetClosedBranch(closedBranchID uint) (*models.ClosedBranch, error) {
	closedBranch, err := branchService.ClosedBranchRepository.GetByID(closedBranchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find closed branch with id %v", closedBranchID)
	}

	return closedBranch, nil
}
