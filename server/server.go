package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
)

type ServiceEnv struct {
	postService    services.PostService
	versionService services.VersionService
	userService    services.UserService
}

type ControllerEnv struct {
	postController    controllers.PostController
	versionController controllers.VersionController
	userController    controllers.UserController
}

func initServiceEnv() ServiceEnv {
	fs := filesystem.InitFilesystem()

	return ServiceEnv{
		postService:    services.PostService{},
		versionService: services.VersionService{Filesystem: fs},
		userService:    services.UserService{},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController:    controllers.PostController{PostService: &serviceEnv.postService},
		versionController: controllers.VersionController{VersionService: &serviceEnv.versionService},
		userController:    controllers.UserController{UserService: &serviceEnv.userService},
	}
}

func Init() {
	_, err := database.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}

	serviceEnv := initServiceEnv()
	controllerEnv := initControllerEnv(&serviceEnv)

	router := SetUpRouter(controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
