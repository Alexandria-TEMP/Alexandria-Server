package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/docs"
)

func SetUpRouter(controllers ControllerEnv) *gin.Engine {
	// Get router
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	err := router.SetTrustedProxies(nil)

	if err != nil {
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

	mergeRequestRouter(v2, controllers)

	discussionRouter(v2, controllers)

	filterRouter := v2.Group("/filter")
	filterRouter.GET("/posts", controllers.filterController.FilterPosts)
	filterRouter.GET("/project-posts", controllers.filterController.FilterProjectPosts)

	tagRouter := v2.Group("/tags")
	tagRouter.GET("/scientific", controllers.tagController.GetScientificTags)
	tagRouter.GET("/completion-status", controllers.tagController.GetCompletionStatusTags)
	tagRouter.GET("/post-type", controllers.tagController.GetPostTypeTags)
	tagRouter.GET("/feedback-preference", controllers.tagController.GetFeedbackPreferenceTags)

	versionRouter(v2, controllers)

	return router
}

func versionRouter(v2 *gin.RouterGroup, controllers ControllerEnv) {
	versionRouter := v2.Group("/versions")
	versionRouter.GET("/:versionID", controllers.versionController.GetVersion)
	versionRouter.POST("", controllers.versionController.CreateVersion)
	versionRouter.GET("/:versionID/render", controllers.versionController.GetRender)
	versionRouter.GET("/:versionID/repository", controllers.versionController.GetRepository)
	versionRouter.GET("/:versionID/tree", controllers.versionController.GetFileTree)
	versionRouter.GET("/:versionID/file/*filepath", controllers.versionController.GetFileFromRepository)
	versionRouter.GET("/:versionID/discussions", controllers.versionController.GetDiscussions)
}

func discussionRouter(v2 *gin.RouterGroup, controllers ControllerEnv) {
	discussionRouter := v2.Group("/discussions")
	discussionRouter.GET("/:discussionID", controllers.discussionController.GetDiscussion)
	discussionRouter.POST("/", controllers.discussionController.CreateDiscussion)
	discussionRouter.DELETE("/:discussionID", controllers.discussionController.DeleteDiscussion)
	discussionRouter.POST("/:discussionID/reports", controllers.discussionController.AddDiscussionReport)
	discussionRouter.GET("/:discussionID/reports", controllers.discussionController.GetDiscussionReports)
	discussionRouter.GET("/reports/:reportID", controllers.discussionController.GetDiscussionReport)
}

func mergeRequestRouter(v2 *gin.RouterGroup, controllers ControllerEnv) {
	mergeRequestRouter := v2.Group("/merge-requests")
	mergeRequestRouter.GET("/:mergeRequestID", controllers.mergeRequestController.GetMergeRequest)
	mergeRequestRouter.POST("/", controllers.mergeRequestController.CreateMergeRequest)
	mergeRequestRouter.PUT("/", controllers.mergeRequestController.UpdateMergeRequest)
	mergeRequestRouter.DELETE("/:mergeRequestID", controllers.mergeRequestController.DeleteMergeRequest)
	mergeRequestRouter.GET("/:mergeRequestID/review-statuses", controllers.mergeRequestController.GetReviewStatus)
	mergeRequestRouter.GET("/reviews/:reviewID", controllers.mergeRequestController.GetReview)
	mergeRequestRouter.POST("/:mergeRequestID/reviews", controllers.mergeRequestController.CreateReview)
	mergeRequestRouter.GET("/:mergeRequestID/can-review/:userID", controllers.mergeRequestController.UserCanReview)
	mergeRequestRouter.GET("/collaborators/:collaboratorID", controllers.mergeRequestController.GetMergeRequestCollaborator)
}

func memberRouter(v2 *gin.RouterGroup, controllers ControllerEnv) {
	memberRouter := v2.Group("/members")
	memberRouter.GET("/:userID", controllers.memberController.GetMember)
	memberRouter.POST("/", controllers.memberController.CreateMember)
	memberRouter.PUT("/", controllers.memberController.UpdateMember)
	memberRouter.DELETE("/:userID", controllers.memberController.DeleteMember)
	memberRouter.GET("/", controllers.memberController.GetAllMembers)
	memberRouter.GET("/:userID/posts", controllers.memberController.GetMemberPosts)
	memberRouter.GET("/:userID/project-posts", controllers.memberController.GetMemberProjectPosts)
	memberRouter.GET("/:userID/merge-requests", controllers.memberController.GetMemberMergeRequests)
	memberRouter.GET("/:userID/discussions", controllers.memberController.GetMemberDiscussions)
	memberRouter.POST("/:userID/saved-posts", controllers.memberController.AddMemberSavedPost)
	memberRouter.POST("/:userID/saved-project-posts", controllers.memberController.AddMemberSavedProjectPost)
	memberRouter.GET("/:userID/saved-posts", controllers.memberController.GetMemberSavedPosts)
	memberRouter.GET("/:userID/saved-project-posts", controllers.memberController.GetMemberSavedProjectPosts)
}

func projectPostRouter(v2 *gin.RouterGroup, controllers ControllerEnv) {
	projectPostRouter := v2.Group("/project-posts")
	projectPostRouter.GET("/:postID", controllers.projectPostController.GetProjectPost)
	projectPostRouter.POST("/", controllers.projectPostController.CreateProjectPost)
	projectPostRouter.PUT("/", controllers.projectPostController.UpdateProjectPost)
	projectPostRouter.DELETE("/:postID", controllers.projectPostController.DeleteProjectPost)
	projectPostRouter.POST("/from-github", controllers.projectPostController.CreateProjectPostFromGithub)
	projectPostRouter.GET("/:postID/all-discussions", controllers.projectPostController.GetProjectPostDiscussions)
	projectPostRouter.GET("/:postID/merge-requests-by-status", controllers.projectPostController.GetProjectPostMRsByStatus)
}

func postRouter(v2 *gin.RouterGroup, controllers ControllerEnv) {
	postRouter := v2.Group("/posts")
	postRouter.GET("/:postID", controllers.postController.GetPost)
	postRouter.POST("/", controllers.postController.CreatePost)
	postRouter.PUT("/", controllers.postController.UpdatePost)
	postRouter.DELETE("/:postID", controllers.postController.DeletePost)
	postRouter.POST("/from-github", controllers.postController.CreatePostFromGithub)
	postRouter.POST("/:postID/reports", controllers.postController.AddPostReport)
	postRouter.GET("/:postID/reports", controllers.postController.GetPostReports)
	postRouter.GET("/reports/:reportID", controllers.postController.GetPostReport)
	postRouter.GET("/collaborators/:collaboratorID", controllers.postController.GetPostCollaborator)
}
