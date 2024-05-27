package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gorm.io/gorm"
)

type RepositoryEnv struct {
	versionRepository database.RepositoryInterface[*models.Version]
}

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

func initRepositoryEnv(db *gorm.DB) RepositoryEnv {
	return RepositoryEnv{
		versionRepository: &database.ModelRepository[*models.Version]{Database: db},
	}
}

func initServiceEnv(repositoryEnv RepositoryEnv, fs *filesystem.Filesystem) ServiceEnv {
	return ServiceEnv{
		postService: services.PostService{},
		versionService: services.VersionService{
			VersionRepository: repositoryEnv.versionRepository,
			Filesystem:        fs,
		},
		userService: services.UserService{},
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
	db, err := database.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.InitFilesystem()

	repositoryEnv := initRepositoryEnv(db)
	serviceEnv := initServiceEnv(repositoryEnv, fs)
	controllerEnv := initControllerEnv(&serviceEnv)

	router := SetUpRouter(controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
