package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
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
	ProjectPostRepository         database.ModelRepositoryInterface[*models.ProjectPost]
	ReviewRepository              database.ModelRepositoryInterface[*models.BranchReview]
	BranchCollaboratorRepository  database.ModelRepositoryInterface[*models.BranchCollaborator]
	DiscussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
	DiscussionRepository          database.ModelRepositoryInterface[*models.Discussion]
	MemberRepository              database.ModelRepositoryInterface[*models.Member]
	Filesystem                    filesystemInterfaces.Filesystem

	RenderService             interfaces.RenderService
	BranchCollaboratorService interfaces.BranchCollaboratorService
}

func (branchService *BranchService) GetBranch(branchID uint) (models.Branch, error) {
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return *branch, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	return *branch, nil
}

func (branchService *BranchService) CreateBranch(branchCreationForm *forms.BranchCreationForm) (models.Branch, error, error) {
	var branch models.Branch

	// verify parent project post exists
	projectPost, err := branchService.ProjectPostRepository.GetByID(branchCreationForm.ProjectPostID)

	if err != nil {
		return branch, fmt.Errorf("no such project post exists"), nil
	}

	// create and save discussion new container
	// we shouldn't have to do this extra, it should be implicit but it isnt...
	discussionContainer := models.DiscussionContainer{}
	if err := branchService.DiscussionContainerRepository.Create(&discussionContainer); err != nil {
		return branch, fmt.Errorf("failed to add discussion container to db"), nil
	}

	// get all collaborators from ids
	collaborators, err := branchService.BranchCollaboratorService.MembersToBranchCollaborators(branchCreationForm.CollaboratingMemberIDs, branchCreationForm.Anonymous)
	if err != nil {
		return branch, err, nil
	}

	// make new branch
	branch = models.Branch{
		UpdatedPostTitle:          branchCreationForm.UpdatedPostTitle,
		UpdatedCompletionStatus:   branchCreationForm.UpdatedCompletionStatus,
		UpdatedScientificFields:   branchCreationForm.UpdatedScientificFields,
		Collaborators:             collaborators,
		DiscussionContainer:       discussionContainer,
		ProjectPostID:             &branchCreationForm.ProjectPostID,
		BranchTitle:               branchCreationForm.BranchTitle,
		RenderStatus:              models.Success,
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}

	// save branch entity to open branches
	projectPost.OpenBranches = append(projectPost.OpenBranches, &branch)

	if _, err := branchService.ProjectPostRepository.Update(projectPost); err != nil {
		return branch, nil, fmt.Errorf("failed to update project post with new branch")
	}

	// set vfs to repository according to the Post of the ProjectPost of the Branch entity
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// create new branch in git repo with branch ID as its name
	if err := branchService.Filesystem.CreateBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		return branch, nil, err
	}

	return branch, nil, nil
}

func (branchService *BranchService) getDiscussionContainerFromIDs(ids []uint) (models.DiscussionContainer, error) {
	discussions := []*models.Discussion{}

	for _, ID := range ids {
		discussion, err := branchService.DiscussionRepository.GetByID(ID)

		if err != nil {
			return models.DiscussionContainer{}, fmt.Errorf("failed to find discussion with id=%v", ID)
		}

		discussions = append(discussions, discussion)
	}

	return models.DiscussionContainer{Discussions: discussions}, nil
}

func (branchService *BranchService) DeleteBranch(branchID uint) error {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
	}

	// checkout repository
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// delete branch
	if err := branchService.Filesystem.DeleteBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return fmt.Errorf("failed to delete branch from vfs with id %v", branchID)
	}

	// delete entity
	if err := branchService.BranchRepository.Delete(branchID); err != nil {
		return fmt.Errorf("failed to find branch with id %v", branchID)
	}

	return nil
}

func (branchService *BranchService) GetReviewStatus(branchID uint) ([]models.BranchReviewDecision, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v", branchID)
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
		return *branchreview, fmt.Errorf("failed to find branch with id %v", reviewID)
	}

	return *branchreview, nil
}

func (branchService *BranchService) CreateReview(form forms.ReviewCreationForm) (models.BranchReview, error) {
	var branchreview models.BranchReview

	// get branch
	branch, err := branchService.BranchRepository.GetByID(form.BranchID)

	if err != nil {
		return branchreview, fmt.Errorf("failed to find branch with id %v", form.BranchID)
	}

	// get member
	member, err := branchService.MemberRepository.GetByID(form.ReviewingMemberID)

	if err != nil {
		return branchreview, fmt.Errorf("failed to find member with id %v", form.ReviewingMemberID)
	}

	// make new branchreview
	branchreview = models.BranchReview{
		BranchID:             form.BranchID,
		Member:               *member,
		BranchReviewDecision: form.BranchReviewDecision,
		Feedback:             form.Feedback,
	}

	// update branch with new branchreview and update branchreview status accordingly
	branch.Reviews = append(branch.Reviews, &branchreview)
	branch.BranchOverallReviewStatus = branchService.UpdateReviewStatus(branch.Reviews)

	// if approved or rejected we close the branch
	if branch.BranchOverallReviewStatus == models.BranchPeerReviewed || branch.BranchOverallReviewStatus == models.BranchRejected {
		if err := branchService.closeBranch(branch); err != nil {
			return branchreview, err
		}

		return branchreview, nil
	}

	// save changes to branch
	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return branchreview, fmt.Errorf("failed to save branch branchreview")
	}

	return branchreview, nil
}

func (branchService *BranchService) closeBranch(branch *models.Branch) error {
	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return fmt.Errorf("failed to get project post")
	}

	// close branch
	closedBranch := &models.ClosedBranch{
		ProjectPostID:        *branch.ProjectPostID,
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

	if branch.UpdatedScientificFields != nil {
		projectPost.Post.ScientificFields = branch.UpdatedScientificFields
	}

	if branch.UpdatedFeedbackPreferences != nil {
		projectPost.ProjectFeedbackPreference = *branch.UpdatedFeedbackPreferences
	}

	// update project post contributors
	branchService.mergeContributors(projectPost, branch.Collaborators)

	// update project post reviewers
	branchService.mergeReviewers(projectPost, branch.Reviews)

	return nil
}

// We add all branch collaborators to the project post as post collaborators with the "reviewer" type, unless they have already been added as such
func (branchService *BranchService) mergeReviewers(projectPost *models.ProjectPost, reviews []*models.BranchReview) {
	// get all member ids which are reviewers present in post collaborators initially
	collaboratorMemberIDs := []uint{}

	for _, c := range projectPost.Post.Collaborators {
		if c.CollaborationType == models.Reviewer {
			collaboratorMemberIDs = append(collaboratorMemberIDs, c.MemberID)
		}
	}

	// add all new post collaborators
	for _, review := range reviews {
		// if the member is already present as a post collaborator, we do not add it again
		if slices.Contains(collaboratorMemberIDs, review.MemberID) {
			continue
		}

		// otherwise we add this post collaborator
		asPostCollaborator := models.PostCollaborator{
			Member:            review.Member,
			PostID:            projectPost.PostID,
			CollaborationType: models.Reviewer,
		}
		projectPost.Post.Collaborators = append(projectPost.Post.Collaborators, &asPostCollaborator)
	}
}

// We add all branch collaborators to the project post as post collaborators with the "contributor" type, unless they have already been added as such
func (branchService *BranchService) mergeContributors(projectPost *models.ProjectPost, branchCollaborators []*models.BranchCollaborator) {
	// get all member ids which are collaborators present in post collaborators initially
	collaboratorMemberIDs := []uint{}

	for _, c := range projectPost.Post.Collaborators {
		if c.CollaborationType == models.Contributor {
			collaboratorMemberIDs = append(collaboratorMemberIDs, c.MemberID)
		}
	}

	// add all new post collaborators
	for _, branchCollaborator := range branchCollaborators {
		// if the member is already present as a post collaborator, we do not add it again
		if slices.Contains(collaboratorMemberIDs, branchCollaborator.MemberID) {
			continue
		}

		// otherwise we add this post collaborator
		asPostCollaborator := models.PostCollaborator{
			Member:            branchCollaborator.Member,
			PostID:            projectPost.PostID,
			CollaborationType: models.Contributor,
		}
		projectPost.Post.Collaborators = append(projectPost.Post.Collaborators, &asPostCollaborator)
	}
}

// UpdateReviewStatus finds the current branchreview status
// If there are 3 approvals, approve the branch
// If there are any rejections, reject the branch
// Otherwise leave open for branchreview
func (branchService *BranchService) UpdateReviewStatus(reviews []*models.BranchReview) models.BranchOverallReviewStatus {
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

func (branchService *BranchService) MemberCanReview(branchID, memberID uint) (bool, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return false, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return false, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
	}

	// get member
	member, err := branchService.MemberRepository.GetByID(memberID)

	if err != nil {
		return false, fmt.Errorf("failed to find member with id %v", memberID)
	}

	// create sets of all tags
	for _, tag := range member.ScientificFields {
		if slices.Contains(branch.UpdatedScientificFields, tag) || slices.Contains(projectPost.Post.ScientificFields, tag) {
			return true, nil
		}
	}

	return false, nil
}

func (branchService *BranchService) GetProject(branchID uint) (string, error) {
	var filePath string

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return filePath, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
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
		return fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
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
		return nil, fmt.Errorf("failed to find branch with id %v", branchID), nil
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return nil, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID), nil
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

func (branchService *BranchService) GetFileFromProject(branchID uint, relFilepath string) (string, error) {
	var absFilepath string

	// validate file path is inside of repository
	if strings.Contains(relFilepath, "..") {
		return absFilepath, fmt.Errorf("file is outside of repository")
	}

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return absFilepath, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(*branch.ProjectPostID)

	if err != nil {
		return absFilepath, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
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

func (branchService *BranchService) GetBranchCollaborator(branchCollaboratorID uint) (*models.BranchCollaborator, error) {
	branchCollaborator, err := branchService.BranchCollaboratorRepository.GetByID(branchCollaboratorID)

	if err != nil {
		return branchCollaborator, fmt.Errorf("failed to get branch collaborator")
	}

	return branchCollaborator, nil
}
