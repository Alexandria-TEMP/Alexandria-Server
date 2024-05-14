package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
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
		postController: controllers.PostController{PostService: &serviceEnv.postService},
	}
}

func Init() {
	db, err := database.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}

	// TODO remove me
	memberRepo := database.ModelRepository[*models.Member]{Database: db}

	err = memberRepo.Create(&models.Member{
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	})

	if err != nil {
		log.Fatal(err)
	}
	// TODO remove until me

	serviceEnv := initServiceEnv()
	controllerEnv := initControllerEnv(serviceEnv)

	router := SetUpRouter(controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
