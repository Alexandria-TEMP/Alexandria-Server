package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
)

type ServiceEnv struct {
	postService services.PostService
	userService services.UserService
}

type ControllerEnv struct {
	postController controllers.PostController
	userController controllers.MemberController
	projectPostController controllers.ProjectPostController
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
		userController: controllers.MemberController{UserService: &serviceEnv.userService},
		projectPostController: controllers.ProjectPostController{},
	}
}

func Init() {
	_, err := database.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}

	serviceEnv := initServiceEnv()
	controllerEnv := initControllerEnv(serviceEnv)

	router := SetUpRouter(controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
