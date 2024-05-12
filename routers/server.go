package routers

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
)

type ServiceEnv struct {
	postService services.PostService
	userService services.UserService
}

type ControllerEnv struct {
	postController controllers.PostController
	userController controllers.UserController
}

func initServiceEnv() ServiceEnv {
	return ServiceEnv{
		postService: services.PostService{},
		userService: services.UserService{},
	}
}

func initControllerEnv(serviceEnv ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController: controllers.PostController{PostService: &serviceEnv.postService},
		userController: controllers.UserController{UserService: &serviceEnv.userService},
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
