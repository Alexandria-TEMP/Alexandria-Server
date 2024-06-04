package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

var (
	cwd    string
	router *gin.Engine

	mockMemberService  *mock_interfaces.MockMemberService
	mockTagService     *mock_interfaces.MockTagService
	memberController   *MemberController
	versionController  *VersionController
	mockVersionService *mock_interfaces.MockVersionService

	responseRecorder *httptest.ResponseRecorder

	exampleMember     models.Member
	exampleMemberDTO  models.MemberDTO
	exampleMemberForm forms.MemberCreationForm
	exampleSTag1      *tags.ScientificFieldTag
	exampleSTag2      *tags.ScientificFieldTag

	examplePendingVersion models.Version
	exampleSuccessVersion models.Version
	exampleFailureVersion models.Version
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = gin.Default()
	gin.SetMode(gin.TestMode)

	router = SetUpRouter()

	examplePendingVersion = models.Version{RenderStatus: models.RenderPending}
	exampleSuccessVersion = models.Version{RenderStatus: models.RenderSuccess}
	exampleFailureVersion = models.Version{RenderStatus: models.RenderFailure}

	tag1 := tags.ScientificFieldTag{
		ScientificField: "Mathematics",
		Subtags:         []*tags.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag1 = &tag1
	tag2 := tags.ScientificFieldTag{
		ScientificField: "",
		Subtags:         []*tags.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag2 = &tag2

	exampleMemberDTO = models.MemberDTO{
		FirstName:             "John",
		LastName:              "Smith",
		Email:                 "john.smith@gmail.com",
		Password:              "password",
		Institution:           "TU Delft",
		ScientificFieldTagIDs: []uint{1, 2},
	}

	exampleMemberForm = forms.MemberCreationForm{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
	}

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
	router.GET("/api/v2/members/:userID", func(c *gin.Context) {
		memberController.GetMember(c)
	})
	router.POST("/api/v2/members", func(c *gin.Context) {
		memberController.CreateMember(c)
	})

	return router
}
