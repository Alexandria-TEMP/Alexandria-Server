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

	versionController  *VersionController
	mockVersionService *mock_interfaces.MockVersionService

	examplePendingVersion models.Version
	exampleSuccessVersion models.Version
	exampleFailureVersion models.Version
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = SetUpRouter()

	examplePendingVersion = models.Version{RenderStatus: models.Pending}
	exampleSuccessVersion = models.Version{RenderStatus: models.Success}
	exampleFailureVersion = models.Version{RenderStatus: models.Failure}

	cwd, _ = os.Getwd()

	os.Exit(m.Run())
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	router.POST("/api/v2/versions", func(c *gin.Context) {
		versionController.CreateVersion(c)
	})
	router.GET("/api/v2/versions/:versionID/render", func(c *gin.Context) {
		versionController.GetRender(c)
	})
	router.GET("/api/v2/versions/:versionID/repository", func(c *gin.Context) {
		versionController.GetRepository(c)
	})
	router.GET("/api/v2/versions/:versionID/tree", func(c *gin.Context) {
		versionController.GetFileTree(c)
	})
	router.GET("/api/v2/versions/:versionID/file/*filepath", func(c *gin.Context) {
		versionController.GetFileFromRepository(c)
	})

	return router
}
