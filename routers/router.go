package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/docs"
)

func SetUpRouter(controllers ControllerEnv) *gin.Engine {
	router := gin.Default()
	err := router.SetTrustedProxies(nil)

	if err != nil {
		return nil
	}

	// Setup swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		postRouter := v1.Group("/post")
		{
			postRouter.GET("/:postID", controllers.postController.GetPost)
		}
	}

	// router.POST("/api/v1/post", controllers.postController.CreatePost)

	return router
}
