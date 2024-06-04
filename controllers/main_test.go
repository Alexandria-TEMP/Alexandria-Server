package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	branchController  *BranchController
	mockBranchService *mock_interfaces.MockBranchService

	examplePendingBranch models.Branch
	exampleSuccessBranch models.Branch
	exampleFailureBranch models.Branch
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = SetUpRouter()

	examplePendingBranch = models.Branch{RenderStatus: models.Pending}
	exampleSuccessBranch = models.Branch{RenderStatus: models.Success}
	exampleFailureBranch = models.Branch{RenderStatus: models.Failure}

	cwd, _ = os.Getwd()

	os.Exit(m.Run())
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	// router.POST("/api/v2/branches", func(c *gin.Context) {
	// 	branchController.CreateVersion(c)
	// })
	// router.GET("/api/v2/versions/:versionID/render", func(c *gin.Context) {
	// 	branchController.GetRender(c)
	// })
	// router.GET("/api/v2/versions/:versionID/repository", func(c *gin.Context) {
	// 	branchController.GetRepository(c)
	// })
	// router.GET("/api/v2/versions/:versionID/tree", func(c *gin.Context) {
	// 	branchController.GetFileTree(c)
	// })
	// router.GET("/api/v2/versions/:versionID/file/*filepath", func(c *gin.Context) {
	// 	branchController.GetFileFromRepository(c)
	// })

	return router
}
