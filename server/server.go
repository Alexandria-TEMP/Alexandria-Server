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
	postRepository               database.ModelRepositoryInterface[*models.Post]
	projectPostRepository        database.ModelRepositoryInterface[*models.ProjectPost]
	memberRepository             database.ModelRepositoryInterface[*models.Member]
	postCollaboratorRepository   database.ModelRepositoryInterface[*models.PostCollaborator]
	branchCollaboratorRepository database.ModelRepositoryInterface[*models.BranchCollaborator]
	tagRepository                database.ModelRepositoryInterface[*tags.ScientificFieldTag]
}

type ServiceEnv struct {
	postService               interfaces.PostService
	projectPostService        interfaces.ProjectPostService
	memberService             interfaces.MemberService
	postCollaboratorService   interfaces.PostCollaboratorService
	branchCollaboratorService interfaces.BranchCollaboratorService
	tagService                interfaces.TagService
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
		postRepository: &database.ModelRepository[*models.Post]{
			Database: db,
		},
		projectPostRepository: &database.ModelRepository[*models.ProjectPost]{
			Database: db,
		},
		memberRepository: &database.ModelRepository[*models.Member]{
			Database: db,
		},
		tagRepository: &database.ModelRepository[*tags.ScientificFieldTag]{
			Database: db,
		},
		postCollaboratorRepository: &database.ModelRepository[*models.PostCollaborator]{
			Database: db,
		},
		branchCollaboratorRepository: &database.ModelRepository[*models.BranchCollaborator]{
			Database: db,
		},
	}
}

func initServiceEnv(repositories *RepositoryEnv, _ *filesystem.Filesystem) *ServiceEnv {
	postCollaboratorService := &services.PostCollaboratorService{
		PostCollaboratorRepository: repositories.postCollaboratorRepository,
		MemberRepository:           repositories.memberRepository,
	}

	branchCollaboratorService := &services.BranchCollaboratorService{
		BranchCollaboratorRepository: repositories.branchCollaboratorRepository,
		MemberRepository:             repositories.memberRepository,
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

	return &ServiceEnv{
		postService:        postService,
		projectPostService: projectPostService,
		memberService: &services.MemberService{
			MemberRepository: repositories.memberRepository,
		},
		tagService: &services.TagService{
			TagRepository: repositories.tagRepository,
		},
		postCollaboratorService:   postCollaboratorService,
		branchCollaboratorService: branchCollaboratorService,
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) *ControllerEnv {
	return &ControllerEnv{
		postController: &controllers.PostController{
			PostService:             serviceEnv.postService,
			PostCollaboratorService: serviceEnv.postCollaboratorService,
		},
		memberController: &controllers.MemberController{
			MemberService: serviceEnv.memberService,
			TagService:    serviceEnv.tagService,
		},
		projectPostController: &controllers.ProjectPostController{
			ProjectPostService: serviceEnv.projectPostService,
		},
		discussionController: &controllers.DiscussionController{},
		filterController: &controllers.FilterController{
			PostService:        serviceEnv.postService,
			ProjectPostService: serviceEnv.projectPostService,
		},
		branchController: &controllers.BranchController{},
		tagController: &controllers.TagController{
			TagService: serviceEnv.tagService,
		},
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
	controllerEnv := initControllerEnv(serviceEnv)

	router := SetUpRouter(controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
