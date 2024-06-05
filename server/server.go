package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gorm.io/gorm"
)

type RepositoryEnv struct {
	postRepository   database.ModelRepositoryInterface[*models.Post]
	memberRepository database.ModelRepositoryInterface[*models.Member]
}

type ServiceEnv struct {
	postService   interfaces.PostService
	memberService interfaces.MemberService
}

type ControllerEnv struct {
	postController        *controllers.PostController
	memberController      *controllers.MemberController
	projectPostController *controllers.ProjectPostController
	discussionController  *controllers.DiscussionController
	filterController      *controllers.FilterController
	branchController      *controllers.BranchController
	tagController         *controllers.TagController
}

func initRepositoryEnv(db *gorm.DB) RepositoryEnv {
	return RepositoryEnv{
		postRepository: &database.ModelRepository[*models.Post]{
			Database: db,
		},
		memberRepository: &database.ModelRepository[*models.Member]{
			Database: db,
		},
	}
}

func initServiceEnv(repositories RepositoryEnv, _ *filesystem.Filesystem) ServiceEnv {
	return ServiceEnv{
		postService: &services.PostService{
			PostRepository:   repositories.postRepository,
			MemberRepository: repositories.memberRepository,
		},
		memberService: &services.MemberService{},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController:        &controllers.PostController{PostService: serviceEnv.postService},
		memberController:      &controllers.MemberController{MemberService: serviceEnv.memberService},
		projectPostController: &controllers.ProjectPostController{},
		discussionController:  &controllers.DiscussionController{},
		filterController:      &controllers.FilterController{},
		branchController:      &controllers.BranchController{},
		tagController:         &controllers.TagController{},
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
