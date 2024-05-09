package routers

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
)

type ServiceEnv struct {
	postService services.PostService
}

type ControllerEnv struct {
	postController controllers.PostController
}

func initServiceEnv() ServiceEnv {
	return ServiceEnv{
		postService: services.PostService{},
	}
}

func initControllerEnv(serviceEnv ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController: controllers.PostController{PostService: serviceEnv.postService},
	}
}

func Init() {
	serviceEnv := initServiceEnv()
	controllerEnv := initControllerEnv(serviceEnv)

	router := SetUpRouter(controllerEnv)
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
