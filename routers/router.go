package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/docs"
)

func SetUpRouter(controllers ControllerEnv) *gin.Engine {
	// Get router
	router := gin.Default()
	err := router.SetTrustedProxies(nil)

	if err != nil {
		return nil
	}

	// Setup swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routing
	v1 := router.Group("/api/v1")

	postRouter := v1.Group("/post")
	postRouter.GET("/:postID", controllers.postController.GetPost)
	postRouter.POST("/", controllers.postController.CreatePost)
	postRouter.PUT("/", controllers.postController.UpdatePost)

	projectPostRouter := v1.Group("/projectPost")
	projectPostRouter.GET("/:postID", controllers.postController.GetProjectPost)
	projectPostRouter.POST("", controllers.postController.CreateProjectPost)

	memberRouter := v1.Group("/member")
	memberRouter.GET("/:userID", controllers.userController.GetMember)
	memberRouter.POST("/", controllers.userController.CreateMember)
	memberRouter.PUT("/", controllers.userController.UpdateMember)

	collaboratorRouter := v1.Group("/collaborator")
	collaboratorRouter.GET("/:userID", controllers.userController.GetCollaborator)
	collaboratorRouter.POST("/", controllers.userController.CreateCollaborator)
	collaboratorRouter.PUT("/", controllers.userController.UpdateCollaborator)

	return router
}
