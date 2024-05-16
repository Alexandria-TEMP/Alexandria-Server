package controllerTests

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

var (
	router          *gin.Engine

	mockPostService *mock_interfaces.MockPostService
	postController  *controllers.PostController

	mockUserService 	*mock_interfaces.MockUserService
	userController  	*controllers.UserController

	responseRecorder       	*httptest.ResponseRecorder
	
	exampleMember          		models.Member
	exampleCollaborator    		models.Collaborator
	exampleMemberForm        	forms.MemberCreationForm
	exampleCollaboratorForm 	forms.CollaboratorCreationForm

	examplePost            models.Post
	exampleProjectPost     models.ProjectPost
	examplePostForm        forms.PostCreationForm
	exampleProjectPostForm forms.ProjectPostCreationForm
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = gin.Default()
	gin.SetMode(gin.TestMode)

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

	// Setup objects
	examplePost = models.Post{ID: 1}
	exampleProjectPost = models.ProjectPost{ID: 2}

	exampleMember = models.Member{}
	exampleCollaborator = models.Collaborator{}

	os.Exit(m.Run())
}