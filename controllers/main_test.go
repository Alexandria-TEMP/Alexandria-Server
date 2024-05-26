package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	postController  *PostController
	mockPostService *mock_interfaces.MockPostService

	examplePostForm        forms.PostCreationForm
	exampleProjectPostForm forms.ProjectPostCreationForm
	examplePost            models.Post
	exampleProjectPost     models.ProjectPost

	userController  *UserController
	mockUserService *mock_interfaces.MockUserService

	exampleMember           models.Member
	exampleCollaborator     models.PostCollaborator
	exampleCollaboratorForm forms.CollaboratorCreationForm
	exampleMemberForm       forms.MemberCreationForm

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

	// Setup objects
	examplePost = models.Post{}
	exampleProjectPost = models.ProjectPost{}

	exampleMember = models.Member{}
	exampleCollaborator = models.PostCollaborator{}

	examplePendingVersion = models.Version{RenderStatus: models.Pending}
	exampleSuccessVersion = models.Version{RenderStatus: models.Success}
	exampleFailureVersion = models.Version{RenderStatus: models.Failure}

	cwd, _ = os.Getwd()

	os.Exit(m.Run())
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	router.GET("/api/v1/post/:postID", func(c *gin.Context) {
		postController.GetPost(c)
	})
	router.POST("/api/v1/post", func(c *gin.Context) {
		postController.CreatePost(c)
	})
	router.PUT("/api/v1/post", func(c *gin.Context) {
		postController.UpdatePost(c)
	})
	router.GET("/api/v1/projectPost/:postID", func(c *gin.Context) {
		postController.GetProjectPost(c)
	})
	router.POST("/api/v1/projectPost", func(c *gin.Context) {
		postController.CreateProjectPost(c)
	})
	router.PUT("/api/v1/projectPost", func(c *gin.Context) {
		postController.UpdateProjectPost(c)
	})
	router.GET("/api/v1/member/:userID", func(c *gin.Context) {
		userController.GetMember(c)
	})
	router.POST("/api/v1/member", func(c *gin.Context) {
		userController.CreateMember(c)
	})
	router.PUT("/api/v1/member", func(c *gin.Context) {
		userController.UpdateMember(c)
	})
	router.GET("/api/v1/collaborator/:userID", func(c *gin.Context) {
		userController.GetCollaborator(c)
	})
	router.POST("/api/v1/collaborator", func(c *gin.Context) {
		userController.CreateCollaborator(c)
	})
	router.PUT("/api/v1/collaborator", func(c *gin.Context) {
		userController.UpdateCollaborator(c)
	})
	router.POST("/api/v1/version/:postID", func(c *gin.Context) {
		versionController.CreateVersion(c)
	})
	router.GET("/api/v1/version/:postID/:versionID/render", func(c *gin.Context) {
		versionController.GetRender(c)
	})
	router.GET("/api/v1/version/:postID/:versionID/repository", func(c *gin.Context) {
		versionController.GetRepository(c)
	})
	router.GET("/api/v1/version/:postID/:versionID/tree", func(c *gin.Context) {
		versionController.GetTreeFromRepository(c)
	})
	router.GET("/api/v1/version/:postID/:versionID/blob/*filepath", func(c *gin.Context) {
		versionController.GetFileFromRepository(c)
	})

	return router
}
