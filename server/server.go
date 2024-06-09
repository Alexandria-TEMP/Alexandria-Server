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
	branchRepository              database.ModelRepositoryInterface[*models.Branch]
	branchCollaboratorRepository  database.ModelRepositoryInterface[*models.BranchCollaborator]
	postRepository                database.ModelRepositoryInterface[*models.Post]
	projectPostRepository         database.ModelRepositoryInterface[*models.ProjectPost]
	reviewRepository              database.ModelRepositoryInterface[*models.Review]
	discussionRepository          database.ModelRepositoryInterface[*models.Discussion]
	discussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
	memberRepository              database.ModelRepositoryInterface[*models.Member]
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

func initRepositoryEnv(db *gorm.DB) *RepositoryEnv {
	return &RepositoryEnv{
		branchRepository:              &database.ModelRepository[*models.Branch]{Database: db},
		branchCollaboratorRepository:  &database.ModelRepository[*models.BranchCollaborator]{Database: db},
		postRepository:                &database.ModelRepository[*models.Post]{Database: db},
		projectPostRepository:         &database.ModelRepository[*models.ProjectPost]{Database: db},
		reviewRepository:              &database.ModelRepository[*models.Review]{Database: db},
		discussionRepository:          &database.ModelRepository[*models.Discussion]{Database: db},
		discussionContainerRepository: &database.ModelRepository[*models.DiscussionContainer]{Database: db},
		memberRepository:              &database.ModelRepository[*models.Member]{Database: db},
	}
}

func initServiceEnv(repositoryEnv *RepositoryEnv, fs *filesystem.Filesystem) ServiceEnv {
	postService := &services.PostService{
		PostRepository:   repositoryEnv.postRepository,
		MemberRepository: repositoryEnv.memberRepository,
	}
	memberService := &services.MemberService{}
	renderService := &services.RenderService{
		BranchRepository:      repositoryEnv.branchRepository,
		ProjectPostRepository: repositoryEnv.projectPostRepository,
		Filesystem:            fs,
	}
	branchService := &services.BranchService{
		BranchRepository:              repositoryEnv.branchRepository,
		ProjectPostRepository:         repositoryEnv.projectPostRepository,
		ReviewRepository:              repositoryEnv.reviewRepository,
		BranchCollaboratorRepository:  repositoryEnv.branchCollaboratorRepository,
		DiscussionContainerRepository: repositoryEnv.discussionContainerRepository,
		DiscussionRepository:          repositoryEnv.discussionRepository,
		MemberRepository:              repositoryEnv.memberRepository,
		Filesystem:                    fs,
		RenderService:                 renderService,
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
