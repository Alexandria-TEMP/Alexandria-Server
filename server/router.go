package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	pagination "github.com/webstradev/gin-pagination"
	docs "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/docs"
)

func SetUpRouter(controllers *ControllerEnv, secret string) *gin.Engine {
	// Get router
	router := gin.Default()
	router.Use(cors.Default())
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}

	// Init middleware
	if err := InitializeMiddleware(controllers.memberController.MemberService, secret); err != nil {
		return nil
	}

	// Setup swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v2"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routing
	v2 := router.Group("/api/v2")

	postRouter(v2, controllers)

	projectPostRouter(v2, controllers)

	memberRouter(v2, controllers)

	branchRouter(v2, controllers)

	filterRouter(v2, controllers)

	tagRouter(v2, controllers)

	discussionRouter(v2, controllers)

	discussionContainerRouter(v2, controllers)

	return router
}

func filterRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	filterRouter := v2.Group("/filter")
	filterRouter.GET("/posts", pagination.Default(), controllers.filterController.FilterPosts)
	filterRouter.GET("/project-posts", pagination.Default(), controllers.filterController.FilterProjectPosts)
}

func tagRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	tagRouter := v2.Group("/tags")
	tagRouter.GET("/scientific", controllers.tagController.GetScientificTags)
	tagRouter.GET("/scientific/:tagID", controllers.tagController.GetScientificFieldTag)
	tagRouter.GET("/completion-status", controllers.tagController.GetCompletionStatusTags)
	tagRouter.GET("/post-type", controllers.tagController.GetPostTypeTags)
	tagRouter.GET("/feedback-preference", controllers.tagController.GetFeedbackPreferenceTags)
}

func discussionRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	discussionRouter := v2.Group("/discussions")
	discussionRouter.GET("/:discussionID", controllers.discussionController.GetDiscussion)
	discussionRouter.POST("/roots", CheckAuth, controllers.discussionController.CreateRootDiscussion)
	discussionRouter.POST("/replies", CheckAuth, controllers.discussionController.CreateReplyDiscussion)
	discussionRouter.DELETE("/:discussionID", CheckAuth, controllers.discussionController.DeleteDiscussion)
	discussionRouter.POST("/:discussionID/reports", CheckAuth, controllers.discussionController.AddDiscussionReport)
	discussionRouter.GET("/:discussionID/reports", controllers.discussionController.GetDiscussionReports)
	discussionRouter.GET("/reports/:reportID", controllers.discussionController.GetDiscussionReport)
}

func branchRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	branchRouter := v2.Group("/branches")
	branchRouter.GET("/:branchID", controllers.branchController.GetBranch)
	branchRouter.POST("", CheckAuth, controllers.branchController.CreateBranch)
	branchRouter.PUT("", CheckAuth, controllers.branchController.UpdateBranch)
	branchRouter.DELETE("/:branchID", CheckAuth, controllers.branchController.DeleteBranch)
	branchRouter.GET("/:branchID/review-statuses", controllers.branchController.GetReviewStatus)
	branchRouter.GET("/reviews/:reviewID", controllers.branchController.GetReview)
	branchRouter.POST("/reviews", CheckAuth, controllers.branchController.CreateReview)
	branchRouter.GET("/:branchID/can-review/:memberID", CheckAuth, controllers.branchController.MemberCanReview)
	branchRouter.GET("/collaborators/:collaboratorID", controllers.branchController.GetBranchCollaborator)
	branchRouter.GET("/:branchID/render", controllers.branchController.GetRender)
	branchRouter.GET("/:branchID/repository", controllers.branchController.GetProject)
	branchRouter.POST("/:branchID/upload", CheckAuth, controllers.branchController.UploadProject)
	branchRouter.GET("/:branchID/tree", controllers.branchController.GetFiletree)
	branchRouter.GET("/:branchID/file/*filepath", controllers.branchController.GetFileFromProject)
	branchRouter.GET("/:branchID/discussions", controllers.branchController.GetDiscussions)
	branchRouter.GET("/closed/:closedBranchID", controllers.branchController.GetClosedBranch)
}

func memberRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	memberRouter := v2.Group("/members")
	memberRouter.GET("/:memberID", controllers.memberController.GetMember)
	memberRouter.POST("", controllers.memberController.CreateMember)
	memberRouter.PUT("", CheckAuth, controllers.memberController.UpdateMember)
	memberRouter.DELETE("/:memberID", CheckAuth, controllers.memberController.DeleteMember)
	memberRouter.GET("", controllers.memberController.GetAllMembers)
	memberRouter.GET("/:memberID/posts", controllers.memberController.GetMemberPosts)
	memberRouter.GET("/:memberID/project-posts", controllers.memberController.GetMemberProjectPosts)
	memberRouter.GET("/:memberID/branches", controllers.memberController.GetMemberBranches)
	memberRouter.GET("/:memberID/discussions", controllers.memberController.GetMemberDiscussions)
	memberRouter.POST("/:memberID/saved-posts", CheckAuth, controllers.memberController.AddMemberSavedPost)
	memberRouter.POST("/:memberID/saved-project-posts", CheckAuth, controllers.memberController.AddMemberSavedProjectPost)
	memberRouter.GET("/:memberID/saved-posts", controllers.memberController.GetMemberSavedPosts)
	memberRouter.GET("/:memberID/saved-project-posts", controllers.memberController.GetMemberSavedProjectPosts)
	memberRouter.GET("/login", controllers.memberController.LoginMember)
	memberRouter.GET("/token", controllers.memberController.RefreshToken)
}

func projectPostRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	projectPostRouter := v2.Group("/project-posts")
	projectPostRouter.GET("/:projectPostID", controllers.projectPostController.GetProjectPost)
	projectPostRouter.POST("", CheckAuth, controllers.projectPostController.CreateProjectPost)
	projectPostRouter.PUT("", CheckAuth, controllers.projectPostController.UpdateProjectPost)
	projectPostRouter.DELETE("/:projectPostID", CheckAuth, controllers.projectPostController.DeleteProjectPost)
	projectPostRouter.POST("/from-github", CheckAuth, controllers.projectPostController.CreateProjectPostFromGithub)
	projectPostRouter.GET("/:projectPostID/all-discussions", controllers.projectPostController.GetProjectPostDiscussions)
	projectPostRouter.GET("/:projectPostID/branches-by-status", controllers.projectPostController.GetProjectPostBranchesByStatus)
}

func postRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	postRouter := v2.Group("/posts")
	postRouter.GET("/:postID", controllers.postController.GetPost)
	postRouter.POST("", CheckAuth, controllers.postController.CreatePost)
	postRouter.PUT("", CheckAuth, controllers.postController.UpdatePost)
	postRouter.DELETE("/:postID", CheckAuth, controllers.postController.DeletePost)
	postRouter.POST("/from-github", CheckAuth, controllers.postController.CreatePostFromGithub)
	postRouter.POST("/:postID/reports", CheckAuth, controllers.postController.AddPostReport)
	postRouter.GET("/:postID/reports", controllers.postController.GetPostReports)
	postRouter.GET("/reports/:reportID", controllers.postController.GetPostReport)
	postRouter.GET("/collaborators/:collaboratorID", controllers.postController.GetPostCollaborator)
	postRouter.POST("/:postID/upload", CheckAuth, controllers.postController.UploadPost)
	postRouter.GET("/:postID/render", controllers.postController.GetMainRender)
	postRouter.GET("/:postID/repository", controllers.postController.GetMainProject)
	postRouter.GET("/:postID/tree", controllers.postController.GetMainFiletree)
	postRouter.GET("/:postID/file/*filepath", controllers.postController.GetMainFileFromProject)
}

func discussionContainerRouter(v2 *gin.RouterGroup, controllers *ControllerEnv) {
	discussionContainerRouter := v2.Group("/discussion-containers")
	discussionContainerRouter.GET("/:discussionContainerID", controllers.discussionContainerController.GetDiscussionContainer)
}
