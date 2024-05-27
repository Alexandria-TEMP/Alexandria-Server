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
	memberService  services.MemberService
}

type ControllerEnv struct {
	postController         controllers.PostController
	memberController       controllers.MemberController
	projectPostController  controllers.ProjectPostController
	discussionController   controllers.DiscussionController
	filterController       controllers.FilterController
	mergeRequestController controllers.MergeRequestController
	tagController          controllers.TagController
	versionController      controllers.VersionController
}

func initServiceEnv() ServiceEnv {
	fs := filesystem.InitFilesystem()

	return ServiceEnv{
		postService:    services.PostService{},
		versionService: services.VersionService{Filesystem: fs},
		memberService:  services.MemberService{},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController:         controllers.PostController{PostService: &serviceEnv.postService},
		memberController:       controllers.MemberController{MemberService: &serviceEnv.memberService},
		projectPostController:  controllers.ProjectPostController{},
		discussionController:   controllers.DiscussionController{},
		filterController:       controllers.FilterController{},
		mergeRequestController: controllers.MergeRequestController{},
		tagController:          controllers.TagController{},
		versionController:      controllers.VersionController{},
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
