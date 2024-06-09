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
	postRepository                database.ModelRepositoryInterface[*models.Post]
	projectPostRepository         database.ModelRepositoryInterface[*models.ProjectPost]
	memberRepository              database.ModelRepositoryInterface[*models.Member]
	discussionRepository          database.ModelRepositoryInterface[*models.Discussion]
	discussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
}

type ServiceEnv struct {
	postService               interfaces.PostService
	projectPostService        interfaces.ProjectPostService
	memberService             interfaces.MemberService
	postCollaboratorService   interfaces.PostCollaboratorService
	branchCollaboratorService interfaces.BranchCollaboratorService
	discussionService         interfaces.DiscussionService
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
		projectPostRepository: &database.ModelRepository[*models.ProjectPost]{
			Database: db,
		},
		memberRepository: &database.ModelRepository[*models.Member]{
			Database: db,
		},
		discussionRepository: &database.ModelRepository[*models.Discussion]{
			Database: db,
		},
		discussionContainerRepository: &database.ModelRepository[*models.DiscussionContainer]{
			Database: db,
		},
	}
}

func initServiceEnv(repositories RepositoryEnv, _ *filesystem.Filesystem) ServiceEnv {
	postCollaboratorService := &services.PostCollaboratorService{
		MemberRepository: repositories.memberRepository,
	}

	branchCollaboratorService := &services.BranchCollaboratorService{
		MemberRepository: repositories.memberRepository,
	}

	postService := &services.PostService{
		PostRepository:          repositories.postRepository,
		MemberRepository:        repositories.memberRepository,
		PostCollaboratorService: postCollaboratorService,
	}

	projectPostService := &services.ProjectPostService{
		ProjectPostRepository:     repositories.projectPostRepository,
		MemberRepository:          repositories.memberRepository,
		PostCollaboratorService:   postCollaboratorService,
		BranchCollaboratorService: branchCollaboratorService,
	}

	discussionService := &services.DiscussionService{
		DiscussionRepository:          repositories.discussionRepository,
		DiscussionContainerRepository: repositories.discussionContainerRepository,
		MemberRepository:              repositories.memberRepository,
	}

	return ServiceEnv{
		postService:               postService,
		projectPostService:        projectPostService,
		memberService:             &services.MemberService{},
		postCollaboratorService:   postCollaboratorService,
		branchCollaboratorService: branchCollaboratorService,
		discussionService:         discussionService,
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController: &controllers.PostController{
			PostService: serviceEnv.postService,
		},
		memberController: &controllers.MemberController{
			MemberService: serviceEnv.memberService,
		},
		projectPostController: &controllers.ProjectPostController{
			ProjectPostService: serviceEnv.projectPostService,
		},
		discussionController: &controllers.DiscussionController{
			DiscussionService: serviceEnv.discussionService,
		},
		filterController: &controllers.FilterController{},
		branchController: &controllers.BranchController{},
		tagController:    &controllers.TagController{},
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
