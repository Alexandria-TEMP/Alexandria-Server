package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-collections/collections/set"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

type BranchService struct {
	BranchRepository       database.RepositoryInterface[*models.Branch]
	ProjectPostRepository  database.RepositoryInterface[*models.ProjectPost]
	BranchReviewRepository database.RepositoryInterface[*models.BranchReview]
	Filesystem             filesystemInterfaces.Filesystem

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

func (branchService *BranchService) CreateBranch(branchCreationForm forms.BranchCreationForm) (models.Branch, error, error) {
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
		ProjectPost:             *projectPost,
		BranchTitle:             branchCreationForm.BranchTitle,
		Anonymous:               branchCreationForm.Anonymous,
	}

	// set vfs to repository according to the Post of the ProjectPost of the Branch entity
	branchService.Filesystem.CheckoutDirectory(branch.ProjectPost.PostID)

	// create new branch in git repo with branch ID as its name
	if err := branchService.Filesystem.CreateBranch(string(branch.ID)); err != nil {
		return branch, nil, err
	}

	// save branch entity
	if err := branchService.BranchRepository.Create(&branch); err != nil {
		return branch, nil, fmt.Errorf("failed to add branch to db")
	}

	return branch, nil, nil
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

func (branchService *BranchService) GetReview(reviewID uint) (models.BranchReview, error) {
	// get branch
	review, err := branchService.BranchReviewRepository.GetByID(reviewID)

	if err != nil {
		return *review, fmt.Errorf("failed to find branch with id %v", reviewID)
	}

	return *review, nil
}

func (branchService *BranchService) CreateReview(form forms.BranchReviewCreationForm) (models.BranchReview, error) {
	var branchReview models.BranchReview

	// get branch
	branch, err := branchService.BranchRepository.GetByID(form.BranchID)

	if err != nil {
		return branchReview, fmt.Errorf("failed to find branch with id %v", form.BranchID)
	}

	// get member
	member, err := branchService.MemberService.GetMember(form.MemberID)

	if err != nil {
		return branchReview, fmt.Errorf("failed to find member with id %v", form.MemberID)
	}

	// make new branch
	branchReview = models.BranchReview{
		BranchID:       form.BranchID,
		Member:         *member,
		BranchDecision: form.BranchDecision,
		Feedback:       form.Feedback,
	}

	// save review
	if err := branchService.BranchReviewRepository.Create(&branchReview); err != nil {
		return branchReview, fmt.Errorf("failed to save branch review")
	}

	// update branch with new review
	branch.Reviews = append(branch.Reviews, &branchReview)

	if _, err := branchService.BranchRepository.Update(branch); err != nil {
		return branchReview, fmt.Errorf("failed to save branch review")
	}

	// TODO: Do we check here if there are 3 positive reviews?

	return branchReview, nil
}

func (branchService *BranchService) MemberCanReview(branchID, memberID uint) (bool, error) {
	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return false, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get member
	member, err := branchService.MemberService.GetMember(memberID)

	if err != nil {
		return false, fmt.Errorf("failed to find member with id %v", memberID)
	}

	// get all tags
	branchTags := set.New(branch.UpdatedScientificFields)
	postTags := set.New(branch.ProjectPost.Post.ScientificFieldTags)
	combinedTags := branchTags.Union(postTags)

	memberTags := set.New(member.ScientificFieldTags)

	// find intersection of tags
	intersection := combinedTags.Intersection(memberTags)

	return intersection.Len() >= 1, nil
}

func (branchService *BranchService) GetProject(branchID uint) (string, error) {
	var filePath string

	// get branch
	branch, err := branchService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(branch.ProjectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(string(branchID)); err != nil {
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

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(branch.ProjectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(string(branchID)); err != nil {
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

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(branch.ProjectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(string(branchID)); err != nil {
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

	// select repository of the parent post
	branchService.Filesystem.CheckoutDirectory(branch.ProjectPost.PostID)

	// checkout specified branch
	if err := branchService.Filesystem.CheckoutBranch(string(branchID)); err != nil {
		return absFilepath, fmt.Errorf("failed to find this git branch, with name %v", branchID)
	}

	absFilepath = filepath.Join(branchService.Filesystem.GetCurrentQuartoDirPath(), relFilepath)

	// Check that file exists, if not return 404
	if exists := utils.FileExists(absFilepath); !exists {
		return "", fmt.Errorf("no such file exists")
	}

	return absFilepath, nil
}
