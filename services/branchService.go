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
	"gorm.io/gorm"
)

type BranchService struct {
	BranchRepository              database.ModelRepositoryInterface[*models.Branch]
	ProjectPostRepository         database.ModelRepositoryInterface[*models.ProjectPost]
	ReviewRepository              database.ModelRepositoryInterface[*models.Review]
	BranchCollaboratorRepository  database.ModelRepositoryInterface[*models.BranchCollaborator]
	DiscussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
	DiscussionRepository          database.ModelRepositoryInterface[*models.Discussion]
	Filesystem                    filesystemInterfaces.Filesystem

	MemberService interfaces.MemberService
	RenderService interfaces.RenderService
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

	// make new branch
	branch = models.Branch{
		NewPostTitle:            branchCreationForm.NewPostTitle,
		UpdatedCompletionStatus: branchCreationForm.UpdatedCompletionStatus,
		UpdatedScientificFields: branchCreationForm.UpdatedScientificFields,
		Collaborators:           branchCreationForm.Collaborators,
		ProjectPostID:           branchCreationForm.ProjectPostID,
		BranchTitle:             branchCreationForm.BranchTitle,
		Anonymous:               branchCreationForm.Anonymous,
	}

	// save branch entity to assign ID
	if err := branchService.BranchRepository.Create(&branch); err != nil {
		return branch, nil, fmt.Errorf("failed to add branch to db")
	}

	// set vfs to repository according to the Post of the ProjectPost of the Branch entity
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// create new branch in git repo with branch ID as its name
	if err := branchService.Filesystem.CreateBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		return branch, nil, err
	}

	return branch, nil, nil
}

func (branchService *BranchService) UpdateBranch(branchDTO models.BranchDTO) (models.Branch, error) {
	var branch models.Branch

	// map collaborator IDs to collaborators
	var collaborators []*models.BranchCollaborator

	for _, ID := range branchDTO.CollaboratorIDs {
		collaborator, err := branchService.BranchCollaboratorRepository.GetByID(ID)

		if err != nil {
			return branch, fmt.Errorf("failed to find branch collaborator with id=%v", ID)
		}

		collaborators = append(collaborators, collaborator)
	}

	// map review IDs to reviews
	var reviews []*models.Review

	for _, ID := range branchDTO.ReviewIDs {
		review, err := branchService.ReviewRepository.GetByID(ID)

		if err != nil {
			return branch, fmt.Errorf("failed to find review with id=%v", ID)
		}

		reviews = append(reviews, review)
	}

	// map discussion IDs to discussion container
	var discussions []*models.Discussion

	for _, ID := range branchDTO.DiscussionIDs {
		discussion, err := branchService.DiscussionRepository.GetByID(ID)

		if err != nil {
			return branch, fmt.Errorf("failed to find discussion with id=%v", ID)
		}

		discussions = append(discussions, discussion)
	}

	discussionContainer := models.DiscussionContainer{Discussions: discussions}

	// check project post exists
	_, err := branchService.ProjectPostRepository.GetByID(branchDTO.ProjectPostID)

	if err != nil {
		return branch, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
	}

	// construct new branch
	branch = models.Branch{
		Model:                   gorm.Model{ID: branchDTO.ID},
		NewPostTitle:            branchDTO.NewPostTitle,
		UpdatedCompletionStatus: branchDTO.UpdatedCompletionStatus,
		UpdatedScientificFields: branchDTO.UpdatedScientificFields,
		Collaborators:           collaborators,
		Reviews:                 reviews,
		DiscussionContainer:     discussionContainer,
		ProjectPostID:           branchDTO.ProjectPostID,
		BranchTitle:             branch.BranchTitle,
		Anonymous:               branchDTO.Anonymous,
		RenderStatus:            branchDTO.RenderStatus,
		ReviewStatus:            branchDTO.ReviewStatus,
	}

	// update entity in DB
	if _, err := branchService.BranchRepository.Update(&branch); err != nil {
		return branch, fmt.Errorf("failed to update old branch with new values in DB")
	}

	return branch, nil
}

func (branchService *BranchService) DeleteBranch(branchID uint) error {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(branch.ProjectPostID)

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

func (branchService *BranchService) GetReviewStatus(branchID uint) ([]models.BranchDecision, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get all decisions
	var decisions []models.BranchDecision
	for _, review := range branch.Reviews {
		decisions = append(decisions, review.BranchDecision)
	}

	return decisions, nil
}

func (branchService *BranchService) GetReview(reviewID uint) (models.Review, error) {
	// get branch
	review, err := branchService.ReviewRepository.GetByID(reviewID)

	if err != nil {
		return *review, fmt.Errorf("failed to find branch with id %v", reviewID)
	}

	return *review, nil
}

func (branchService *BranchService) CreateReview(form forms.ReviewCreationForm) (models.Review, error) {
	var review models.Review

	// get branch
	branch, err := branchService.BranchRepository.GetByID(form.BranchID)

	if err != nil {
		return review, fmt.Errorf("failed to find branch with id %v", form.BranchID)
	}

	// get member
	member, err := branchService.MemberService.GetMember(form.MemberID)

	if err != nil {
		return review, fmt.Errorf("failed to find member with id %v", form.MemberID)
	}

	// make new branch
	review = models.Review{
		BranchID:       form.BranchID,
		Member:         *member,
		BranchDecision: form.BranchDecision,
		Feedback:       form.Feedback,
	}

	// save review
	if err := branchService.ReviewRepository.Create(&review); err != nil {
		return review, fmt.Errorf("failed to save branch review")
	}

	// update branch with new review
	branch.Reviews = append(branch.Reviews, &review)

	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return review, fmt.Errorf("failed to save branch review")
	}

	// TODO: Do we check here if there are 3 positive reviews?

	return review, nil
}

func (branchService *BranchService) MemberCanReview(branchID, memberID uint) (bool, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return false, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(branch.ProjectPostID)

	if err != nil {
		return false, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
	}

	// get member
	member, err := branchService.MemberService.GetMember(memberID)

	if err != nil {
		return false, fmt.Errorf("failed to find member with id %v", memberID)
	}

	// create sets of all tags
	for _, tag := range member.ScientificFieldTags {
		if slices.Contains(branch.UpdatedScientificFields, tag) || slices.Contains(projectPost.Post.ScientificFieldTags, tag) {
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
	projectPost, err := branchService.ProjectPostRepository.GetByID(branch.ProjectPostID)

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
	projectPost, err := branchService.ProjectPostRepository.GetByID(branch.ProjectPostID)

	if err != nil {
		return fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return fmt.Errorf("failed to find this git branch, with name %v", branchID)
	}

	// clean directory to remove all files
	if err := branchService.Filesystem.CleanDir(); err != nil {
		return fmt.Errorf("failed to remove all old files")
	}

	// save zipped project
	if err := branchService.Filesystem.SaveZipFile(c, file); err != nil {
		// it fails so we set render status to failed and reset the branch
		branch.RenderStatus = models.Failure
		_, _ = branchService.BranchRepository.Update(branch)
		_ = branchService.Filesystem.Reset()

		return fmt.Errorf("failed to remove all old files")
	}

	go branchService.RenderService.Render(branch)

	return nil
}

func (branchService *BranchService) GetFiletree(branchID uint) (map[string]int64, error, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return nil, fmt.Errorf("failed to find branch with id %v", branchID), nil
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(branch.ProjectPostID)

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

	if strings.Contains(relFilepath, "..") {
		return absFilepath, fmt.Errorf("file is outside of repository")
	}

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return absFilepath, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := branchService.ProjectPostRepository.GetByID(branch.ProjectPostID)

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
