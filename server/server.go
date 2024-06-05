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
	branchRepository      database.ModelRepositoryInterface[*models.Branch]
	postRepository        database.ModelRepositoryInterface[*models.Post]
	projectPostRepository database.ModelRepositoryInterface[*models.ProjectPost]
	reviewRepository      database.ModelRepositoryInterface[*models.Review]
}

type ServiceEnv struct {
	postService   interfaces.PostService
	memberService interfaces.MemberService
	branchService interfaces.BranchService
	renderService interfaces.RenderService
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
		branchRepository: &database.ModelRepository[*models.Branch]{Database: db},
		postRepository:   &database.ModelRepository[*models.Post]{Database: db},
	}
}

func initServiceEnv(repositoryEnv RepositoryEnv, fs *filesystem.Filesystem) ServiceEnv {
	postService := &services.PostService{}
	memberService := &services.MemberService{}
	renderService := &services.RenderService{
		BranchRepository:      repositoryEnv.branchRepository,
		ProjectPostRepository: repositoryEnv.projectPostRepository,
		Filesystem:            fs,
	}
	branchService := &services.BranchService{
		BranchRepository:      repositoryEnv.branchRepository,
		ProjectPostRepository: repositoryEnv.projectPostRepository,
		ReviewRepository:      repositoryEnv.reviewRepository,
		Filesystem:            fs,
		RenderService:         renderService,
		MemberService:         memberService,
	}

	return ServiceEnv{
		postService:   postService,
		memberService: memberService,
		renderService: renderService,
		branchService: branchService,
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController:        &controllers.PostController{PostService: serviceEnv.postService},
		memberController:      &controllers.MemberController{MemberService: serviceEnv.memberService},
		projectPostController: &controllers.ProjectPostController{},
		discussionController:  &controllers.DiscussionController{},
		filterController:      &controllers.FilterController{},
		branchController: &controllers.BranchController{
			BranchService: serviceEnv.branchService,
			RenderService: serviceEnv.renderService,
		},
		tagController: &controllers.TagController{},
	}
}

func Init() {
	db, err := database.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.NewFilesystem()

	repositoryEnv := initRepositoryEnv(db)
	serviceEnv := initServiceEnv(repositoryEnv, fs)
	controllerEnv := initControllerEnv(&serviceEnv)

	router := SetUpRouter(controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
