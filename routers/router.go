package routers

import (
	"github.com/gin-gonic/gin"
)

func SetUpRouter(controllers ControllerEnv) *gin.Engine {
	router := gin.Default()
	err := router.SetTrustedProxies(nil)

	if err != nil {
		return nil
	}

	router.GET("/post/:postID", controllers.postController.GetPost)

	return router
}
