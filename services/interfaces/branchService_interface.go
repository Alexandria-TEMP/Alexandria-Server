package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./branchService_interface.go -destination=../../mocks/branchService_mock.go

type BranchService interface {
	// GetBranch gets an existing branch from the DB
	GetBranch(branchID uint) (models.Branch, error)

	// CreateBranch creates a new branch from a creation form.
	// It assumes that a repository has already been created for this post.
	// Error 1 404
	// Error 2 500
	CreateBranch(branchCreationForm *forms.BranchCreationForm, member *models.Member) (models.Branch, error, error)

	// DeleteBranch deletes an existing branch entity, as well as the branch in the vfs.
	DeleteBranch(branchID uint) error

	// GetAllBranchReviewStatuses gets the decisions for all reviews of a branch, given its ID,
	GetAllBranchReviewStatuses(branchID uint) ([]models.BranchReviewDecision, error)

	// GetReview gets an existing branchreview from the DB
	GetReview(reviewID uint) (models.BranchReview, error)

	// CreateReview creates a new branchreview and adds it to the branch.
	CreateReview(reviewCreationForm forms.ReviewCreationForm, reviewingMember *models.Member) (models.BranchReview, error)

	// MemberCanReview checks whether a user is elligible to branchreview a branch, dpending on whether there is an overlap of the scientific fields.
	MemberCanReview(branchID, memberID uint) (bool, error)

	// GetProjectFile returns filepath of zipped repository.
	// Error is for status 404.
	GetProject(branchID uint) (string, error)

	// UploadProject saves a zipped quarto project to its branch and starts the render pipeline.
	// It renders the project in a goroutine.
	UploadProject(c *gin.Context, file *multipart.FileHeader, branchID uint) error

	// GetFiletree returns a map of all filepaths in a quarto project and their size in bytes
	// Error 1 is for status 404.
	// Error 2 is for status 500.
	GetFiletree(branchID uint) (map[string]int64, error, error)

	// GetFileFromRepository returns absolute filepath of file.
	// Error is for status 404.
	GetFileFromProject(branchID uint, relFilepath string) (string, error)

	// GetBranchProjectPost returns a deeply preloaded project post for a branch.
	GetBranchProjectPost(branch *models.Branch) (*models.ProjectPost, error)

	GetClosedBranch(closedBranchID uint) (*models.ClosedBranch, error)
}
