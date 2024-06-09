package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	branchController  BranchController
	mockBranchService *mocks.MockBranchService
	mockRenderService *mocks.MockRenderService

	exampleBranch       models.Branch
	exampleReview       models.Review
	exampleCollaborator models.BranchCollaborator
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = SetUpRouter()

	cwd, _ = os.Getwd()

	os.Exit(m.Run())
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	branchRouter := router.Group("/api/v2/branches")
	branchRouter.GET("/:branchID", branchController.GetBranch)
	branchRouter.POST("", branchController.CreateBranch)
	branchRouter.PUT("", branchController.UpdateBranch)
	branchRouter.DELETE("/:branchID", branchController.DeleteBranch)
	branchRouter.GET("/:branchID/review-statuses", branchController.GetReviewStatus)
	branchRouter.GET("/reviews/:reviewID", branchController.GetReview)
	branchRouter.POST("/reviews", branchController.CreateReview)
	branchRouter.GET("/:branchID/can-review/:memberID", branchController.MemberCanReview)
	branchRouter.GET("/collaborators/:collaboratorID", branchController.GetBranchCollaborator)
	branchRouter.GET("/:branchID/render", branchController.GetRender)
	branchRouter.GET("/:branchID/repository", branchController.GetProject)
	branchRouter.POST("/:branchID", branchController.UploadProject)
	branchRouter.GET("/:branchID/tree", branchController.GetFiletree)
	branchRouter.GET("/:branchID/file/*filepath", branchController.GetFileFromProject)
	branchRouter.GET("/:branchID/discussions", branchController.GetDiscussions)

	return router
}
