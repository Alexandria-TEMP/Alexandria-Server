package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	pagination "github.com/webstradev/gin-pagination"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	branchController              BranchController
	postController                PostController
	projectPostController         ProjectPostController
	memberController              MemberController
	tagController                 TagController
	discussionContainerController DiscussionContainerController
	filterController              FilterController
	discussionController          DiscussionController

	mockBranchService                      *mocks.MockBranchService
	mockRenderService                      *mocks.MockRenderService
	mockBranchCollaboratorService          *mocks.MockBranchCollaboratorService
	mockMemberService                      *mocks.MockMemberService
	mockTagService                         *mocks.MockTagService
	mockScientificFieldTagContainerService *mocks.MockScientificFieldTagContainerService
	mockPostCollaboratorService            *mocks.MockPostCollaboratorService
	mockPostService                        *mocks.MockPostService
	mockDiscussionContainerService         *mocks.MockDiscussionContainerService

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

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	v2 := router.Group("/api/v2")

	postRouter(v2, &postController)
	projectPostRouter(v2, &projectPostController)
	memberRouter(v2, &memberController)
	branchRouter(v2, &branchController)
	filterRouter(v2, &filterController)
	tagRouter(v2, &tagController)
	discussionRouter(v2, &discussionController)
	discussionContainerRouter(v2, &discussionContainerController)

	return router
}

func filterRouter(v2 *gin.RouterGroup, controller *FilterController) {
	filterRouter := v2.Group("/filter")
	filterRouter.GET("/posts", pagination.Default(), controller.FilterPosts)
	filterRouter.GET("/project-posts", pagination.Default(), controller.FilterProjectPosts)
}

func tagRouter(v2 *gin.RouterGroup, controller *TagController) {
	tagRouter := v2.Group("/tags")
	tagRouter.GET("/scientific", controller.GetScientificTags)
	tagRouter.GET("/scientific/:tagID", controller.GetScientificFieldTag)
	tagRouter.GET("/scientific/containers/:containerID", controller.GetScientificFieldTagContainer)
	tagRouter.GET("/completion-status", controller.GetCompletionStatusTags)
	tagRouter.GET("/post-type", controller.GetPostTypeTags)
	tagRouter.GET("/feedback-preference", controller.GetFeedbackPreferenceTags)
}

func discussionRouter(v2 *gin.RouterGroup, controller *DiscussionController) {
	discussionRouter := v2.Group("/discussions")
	discussionRouter.GET("/:discussionID", controller.GetDiscussion)
	discussionRouter.POST("/roots", controller.CreateRootDiscussion)
	discussionRouter.POST("/replies", controller.CreateReplyDiscussion)
	discussionRouter.DELETE("/:discussionID", controller.DeleteDiscussion)
	discussionRouter.POST("/:discussionID/reports", controller.AddDiscussionReport)
	discussionRouter.GET("/:discussionID/reports", controller.GetDiscussionReports)
	discussionRouter.GET("/reports/:reportID", controller.GetDiscussionReport)
}

func branchRouter(v2 *gin.RouterGroup, controller *BranchController) {
	branchRouter := v2.Group("/branches")
	branchRouter.GET("/:branchID", controller.GetBranch)
	branchRouter.POST("", controller.CreateBranch)
	branchRouter.DELETE("/:branchID", controller.DeleteBranch)
	branchRouter.GET("/:branchID/review-statuses", controller.GetAllBranchReviewStatuses)
	branchRouter.GET("/reviews/:reviewID", controller.GetReview)
	branchRouter.POST("/reviews", controller.CreateReview)
	branchRouter.GET("/:branchID/can-review/:memberID", controller.MemberCanReview)
	branchRouter.GET("/collaborators/:collaboratorID", controller.GetBranchCollaborator)
	branchRouter.GET("/collaborators/all/:branchID", controller.GetAllBranchCollaborators)
	branchRouter.GET("/:branchID/render", controller.GetRender)
	branchRouter.GET("/:branchID/repository", controller.GetProject)
	branchRouter.POST("/:branchID/upload", controller.UploadProject)
	branchRouter.GET("/:branchID/tree", controller.GetFiletree)
	branchRouter.GET("/:branchID/file/*filepath", controller.GetFileFromProject)
	branchRouter.GET("/:branchID/discussions", controller.GetDiscussions)
	branchRouter.GET("/closed/:closedBranchID", controller.GetClosedBranch)
}

func memberRouter(v2 *gin.RouterGroup, controller *MemberController) {
	memberRouter := v2.Group("/members")
	memberRouter.GET("/:memberID", controller.GetMember)
	memberRouter.POST("", controller.CreateMember)
	memberRouter.DELETE("/:memberID", controller.DeleteMember)
	memberRouter.GET("", controller.GetAllMembers)
	memberRouter.GET("/:memberID/posts", controller.GetMemberPosts)
	memberRouter.GET("/:memberID/project-posts", controller.GetMemberProjectPosts)
	memberRouter.GET("/:memberID/branches", controller.GetMemberBranches)
	memberRouter.GET("/:memberID/discussions", controller.GetMemberDiscussions)
	memberRouter.POST("/:memberID/saved-posts", controller.AddMemberSavedPost)
	memberRouter.POST("/:memberID/saved-project-posts", controller.AddMemberSavedProjectPost)
	memberRouter.GET("/:memberID/saved-posts", controller.GetMemberSavedPosts)
	memberRouter.GET("/:memberID/saved-project-posts", controller.GetMemberSavedProjectPosts)
}

func projectPostRouter(v2 *gin.RouterGroup, controller *ProjectPostController) {
	projectPostRouter := v2.Group("/project-posts")
	projectPostRouter.GET("/:projectPostID", controller.GetProjectPost)
	projectPostRouter.POST("", controller.CreateProjectPost)
	projectPostRouter.DELETE("/:projectPostID", controller.DeleteProjectPost)
	projectPostRouter.POST("/from-github", controller.CreateProjectPostFromGithub)
	projectPostRouter.GET("/:projectPostID/all-discussion-containers", controller.GetProjectPostDiscussionContainers)
	projectPostRouter.GET("/:projectPostID/branches-by-status", controller.GetProjectPostBranchesByStatus)
}

func postRouter(v2 *gin.RouterGroup, controller *PostController) {
	postRouter := v2.Group("/posts")
	postRouter.GET("/:postID", controller.GetPost)
	postRouter.POST("", controller.CreatePost)
	postRouter.DELETE("/:postID", controller.DeletePost)
	postRouter.POST("/from-github", controller.CreatePostFromGithub)
	postRouter.POST("/:postID/reports", controller.AddPostReport)
	postRouter.GET("/:postID/reports", controller.GetPostReports)
	postRouter.GET("/reports/:reportID", controller.GetPostReport)
	postRouter.GET("/collaborators/:collaboratorID", controller.GetPostCollaborator)
	postRouter.GET("/collaborators/all/:postID", controller.GetAllPostCollaborators)
	postRouter.POST("/:postID/upload", controller.UploadPost)
	postRouter.GET("/:postID/render", controller.GetMainRender)
	postRouter.GET("/:postID/repository", controller.GetMainProject)
	postRouter.GET("/:postID/tree", controller.GetMainFiletree)
	postRouter.GET("/:postID/file/*filepath", controller.GetMainFileFromProject)
	postRouter.GET("/:postID/project-post", controller.GetProjectPostIfExists)
}

func discussionContainerRouter(v2 *gin.RouterGroup, controller *DiscussionContainerController) {
	discussionContainerRouter := v2.Group("/discussion-containers")
	discussionContainerRouter.GET("/:discussionContainerID", controller.GetDiscussionContainer)
}
