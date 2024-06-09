package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gorm.io/gorm"
)

type RepositoryEnv struct {
	postRepository   database.ModelRepositoryInterface[*models.Post]
	memberRepository database.ModelRepositoryInterface[*models.Member]
	tagRepository    database.ModelRepositoryInterface[*tags.ScientificFieldTag]
}

type ServiceEnv struct {
	postService   interfaces.PostService
	memberService interfaces.MemberService
	tagService    interfaces.TagService
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

func initRepositoryEnv(_ *gorm.DB) RepositoryEnv {
	return RepositoryEnv{}
}

func initServiceEnv(repositoryEnv RepositoryEnv, _ *filesystem.Filesystem) ServiceEnv {
	return ServiceEnv{
		postService: &services.PostService{
			PostRepository:   repositoryEnv.postRepository,
			MemberRepository: repositoryEnv.memberRepository,
		},
		memberService: &services.MemberService{
			MemberRepository: repositoryEnv.memberRepository,
		},
		tagService: &services.TagService{
			TagRepository: repositoryEnv.tagRepository,
		},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController: &controllers.PostController{PostService: serviceEnv.postService},
		memberController: &controllers.MemberController{
			MemberService: serviceEnv.memberService,
			TagService:    serviceEnv.tagService,
		},
		projectPostController: &controllers.ProjectPostController{},
		discussionController:  &controllers.DiscussionController{},
		filterController:      &controllers.FilterController{},
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
