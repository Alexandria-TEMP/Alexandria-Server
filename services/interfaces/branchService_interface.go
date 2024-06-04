package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchService interface {
	// GetBranch gets an existing branch from the DB
	GetBranch(branchID uint) (models.Branch, error)

	// CreateBranch creates a new branch from a creation form.
	// It assumes that a repository has already been created for this post.
	// Error 1 404
	// Error 2 500
	CreateBranch(branchCreationForm forms.BranchCreationForm) (models.Branch, error, error)

	// GetReviewStatus gets the decisions for all reviews of a branch, given its ID,
	GetReviewStatus(branchID uint) ([]models.BranchDecision, error)

	// GetReview gets an existing review from the DB
	GetReview(reviewID uint) (models.Review, error)

	// CreateReview creates a new review and adds it to the branch.
	CreateReview(branchID uint, reportCreationForm forms.ReviewCreationForm) (models.Review, error)

	// UserCanReview checks whether a user is elligible to review a branch, dpending on whether there is an overlap of the scientific fields.
	UserCanReview(branchID, userID uint) (bool, error)

	// UploadProject saves a zipper quarto project to its branch and sets the branch to pending.
	// It the renders the project in a goroutine.
	UploadProject(c *gin.Context, file *multipart.FileHeader, branchID uint) error

	// GetProjectFile returns filepath of zipped repository.
	// Error is for status 404.
	GetProject(branchID uint) (string, error)

	// GetFiletree returns a map of all filepaths in a quarto project and their size in bytes
	// Error 1 is for status 404.
	// Error 2 is for status 500.
	GetFiletree(branchID uint) (map[string]int64, error, error)

	// GetFileFromRepository returns absolute filepath of file.
	// Error is for status 404.
	GetFileFromProject(branchID uint, relFilepath string) (string, error)
}
