package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	branchController BranchController
	postController   PostController
	memberController *MemberController
	tagController    TagController

	mockBranchService                      *mocks.MockBranchService
	mockRenderService                      *mocks.MockRenderService
	mockBranchCollaboratorService          *mocks.MockBranchCollaboratorService
	mockMemberService                      *mocks.MockMemberService
	mockTagService                         *mocks.MockTagService
	mockScientificFieldTagContainerService *mocks.MockScientificFieldTagContainerService
	mockPostCollaboratorService            *mocks.MockPostCollaboratorService
	mockPostService                        *mocks.MockPostService

	exampleBranch       models.Branch
	exampleReview       models.BranchReview
	exampleCollaborator models.BranchCollaborator
	exampleMember       models.Member
	exampleMemberDTO    models.MemberDTO
	exampleMemberForm   forms.MemberCreationForm
	exampleSTag1        *models.ScientificFieldTag
	exampleSTag2        *models.ScientificFieldTag
	exampleSTag1DTO     models.ScientificFieldTagDTO
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	exampleSTag1 = &models.ScientificFieldTag{
		ScientificField: "Mathematics",
		Subtags:         []*models.ScientificFieldTag{},
	}
	exampleSTag2 = &models.ScientificFieldTag{
		ScientificField: "Computers",
		Subtags:         []*models.ScientificFieldTag{},
	}
	exampleSTag1DTO = models.ScientificFieldTagDTO{
		ScientificField: "Mathematics",
		SubtagIDs:       []uint{},
	}
	exampleMember = models.Member{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
		ScientificFieldTagContainer: models.ScientificFieldTagContainer{
			ScientificFieldTags: []*models.ScientificFieldTag{},
		},
	}
	exampleMemberDTO = models.MemberDTO{
		FirstName:                     "John",
		LastName:                      "Smith",
		Email:                         "john.smith@gmail.com",
		Password:                      "password",
		Institution:                   "TU Delft",
		ScientificFieldTagContainerID: 0,
	}

	exampleMemberForm = forms.MemberCreationForm{
		FirstName:             "John",
		LastName:              "Smith",
		Email:                 "john.smith@gmail.com",
		Password:              "password",
		Institution:           "TU Delft",
		ScientificFieldTagIDs: []uint{},
	}

	// Setup test router, to test controller endpoints through http
	router = SetUpRouter()

	cwd, _ = os.Getwd()

	os.Exit(m.Run())
}

// TODO this duplicates a LOT of server logic and so is a pain to maintain...
// TODO could we call the actual server routing function (in router.go) instead?
func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	branchRouter := router.Group("/api/v2/branches")
	branchRouter.GET("/:branchID", branchController.GetBranch)
	branchRouter.POST("", branchController.CreateBranch)
	branchRouter.PUT("", branchController.UpdateBranch)
	branchRouter.DELETE("/:branchID", branchController.DeleteBranch)
	branchRouter.GET("/:branchID/branchreview-statuses", branchController.GetReviewStatus)
	branchRouter.GET("/reviews/:reviewID", branchController.GetReview)
	branchRouter.POST("/reviews", branchController.CreateReview)
	branchRouter.GET("/:branchID/can-branchreview/:memberID", branchController.MemberCanReview)
	branchRouter.GET("/collaborators/:collaboratorID", branchController.GetBranchCollaborator)
	branchRouter.GET("/collaborators/all/:branchID", branchController.GetAllBranchCollaborators)
	branchRouter.GET("/:branchID/render", branchController.GetRender)
	branchRouter.GET("/:branchID/repository", branchController.GetProject)
	branchRouter.POST("/:branchID", branchController.UploadProject)
	branchRouter.GET("/:branchID/tree", branchController.GetFiletree)
	branchRouter.GET("/:branchID/file/*filepath", branchController.GetFileFromProject)
	branchRouter.GET("/:branchID/discussions", branchController.GetDiscussions)

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
	router.GET("/api/v2/tags/scientific/:tagID", func(c *gin.Context) {
		tagController.GetScientificFieldTag(c)
	})
	router.GET("/api/v2/tags/scientific/containers/:containerID", tagController.GetScientificFieldTagContainer)

	postRouter := router.Group("/api/v2/posts")
	postRouter.GET("/collaborators/all/:postID", postController.GetAllPostCollaborators)

	return router
}
