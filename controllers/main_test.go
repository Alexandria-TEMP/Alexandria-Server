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

	mockMemberService *mock_interfaces.MockMemberService
	mockTagService    *mock_interfaces.MockTagService
	memberController  *MemberController
	tagController     *TagController

	responseRecorder *httptest.ResponseRecorder

	exampleMember     models.Member
	exampleMemberDTO  models.MemberDTO
	exampleMemberForm forms.MemberCreationForm
	exampleSTag1      *tags.ScientificFieldTag
	exampleSTag2      *tags.ScientificFieldTag
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = gin.Default()
	gin.SetMode(gin.TestMode)

	router = SetUpRouter()

	tag1 := tags.ScientificFieldTag{
		ScientificField: "Mathematics",
		Subtags:         []*tags.ScientificFieldTag{},
	}
	exampleSTag1 = &tag1
	tag2 := tags.ScientificFieldTag{
		ScientificField: "Computers",
		Subtags:         []*tags.ScientificFieldTag{},
	}
	exampleSTag2 = &tag2

	exampleMember = models.Member{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
	}
	exampleMemberDTO = models.MemberDTO{
		FirstName:             "John",
		LastName:              "Smith",
		Email:                 "john.smith@gmail.com",
		Password:              "password",
		Institution:           "TU Delft",
		ScientificFieldTagIDs: []uint{},
	}

	exampleMemberForm = forms.MemberCreationForm{
		FirstName:             "John",
		LastName:              "Smith",
		Email:                 "john.smith@gmail.com",
		Password:              "password",
		Institution:           "TU Delft",
		ScientificFieldTagIDs: []uint{},
	}

	cwd, _ = os.Getwd()
	os.Exit(m.Run())
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	router.GET("/api/v2/members/:memberID", func(c *gin.Context) {
		memberController.GetMember(c)
	})
	router.POST("/api/v2/members", func(c *gin.Context) {
		memberController.CreateMember(c)
	})
	router.PUT("/api/v2/members", func(c *gin.Context) {
		memberController.UpdateMember(c)
	})
	router.DELETE("/api/v2/members/:memberID", func(c *gin.Context) {
		memberController.DeleteMember(c)
	})
	router.GET("/api/v2/members", func(c *gin.Context) {
		memberController.GetAllMembers(c)
	})
	router.GET("/api/v2/tags/scientific", func(c *gin.Context) {
		tagController.GetScientificTags(c)
	})
	return router
}
