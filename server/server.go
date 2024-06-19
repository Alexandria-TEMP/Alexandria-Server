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
	branchRepository                      database.ModelRepositoryInterface[*models.Branch]
	closedBranchRepository                database.ModelRepositoryInterface[*models.ClosedBranch]
	branchCollaboratorRepository          database.ModelRepositoryInterface[*models.BranchCollaborator]
	postCollaboratorRepository            database.ModelRepositoryInterface[*models.PostCollaborator]
	postRepository                        database.ModelRepositoryInterface[*models.Post]
	projectPostRepository                 database.ModelRepositoryInterface[*models.ProjectPost]
	reviewRepository                      database.ModelRepositoryInterface[*models.BranchReview]
	discussionRepository                  database.ModelRepositoryInterface[*models.Discussion]
	discussionContainerRepository         database.ModelRepositoryInterface[*models.DiscussionContainer]
	memberRepository                      database.ModelRepositoryInterface[*models.Member]
	scientificFieldTagRepository          database.ModelRepositoryInterface[*models.ScientificFieldTag]
	scientificFieldTagContainerRepository database.ModelRepositoryInterface[*models.ScientificFieldTagContainer]
}

type ServiceEnv struct {
	postService                        interfaces.PostService
	memberService                      interfaces.MemberService
	branchService                      interfaces.BranchService
	renderService                      interfaces.RenderService
	projectPostService                 interfaces.ProjectPostService
	postCollaboratorService            interfaces.PostCollaboratorService
	branchCollaboratorService          interfaces.BranchCollaboratorService
	discussionService                  interfaces.DiscussionService
	discussionContainerService         interfaces.DiscussionContainerService
	tagService                         interfaces.TagService
	scientificFieldTagContainerService interfaces.ScientificFieldTagContainerService
}

type ControllerEnv struct {
	postController                *controllers.PostController
	memberController              *controllers.MemberController
	projectPostController         *controllers.ProjectPostController
	discussionController          *controllers.DiscussionController
	discussionContainerController *controllers.DiscussionContainerController
	filterController              *controllers.FilterController
	branchController              *controllers.BranchController
	tagController                 *controllers.TagController
}

func initRepositoryEnv(db *gorm.DB) *RepositoryEnv {
	return &RepositoryEnv{
		branchRepository:                      &database.ModelRepository[*models.Branch]{Database: db},
		closedBranchRepository:                &database.ModelRepository[*models.ClosedBranch]{Database: db},
		branchCollaboratorRepository:          &database.ModelRepository[*models.BranchCollaborator]{Database: db},
		postCollaboratorRepository:            &database.ModelRepository[*models.PostCollaborator]{Database: db},
		postRepository:                        &database.ModelRepository[*models.Post]{Database: db},
		projectPostRepository:                 &database.ModelRepository[*models.ProjectPost]{Database: db},
		reviewRepository:                      &database.ModelRepository[*models.BranchReview]{Database: db},
		discussionRepository:                  &database.ModelRepository[*models.Discussion]{Database: db},
		discussionContainerRepository:         &database.ModelRepository[*models.DiscussionContainer]{Database: db},
		memberRepository:                      &database.ModelRepository[*models.Member]{Database: db},
		scientificFieldTagRepository:          &database.ModelRepository[*models.ScientificFieldTag]{Database: db},
		scientificFieldTagContainerRepository: &database.ModelRepository[*models.ScientificFieldTagContainer]{Database: db},
	}
}

func initServiceEnv(repositoryEnv *RepositoryEnv, fs *filesystem.Filesystem) ServiceEnv {
	tagService := &services.TagService{
		TagRepository: repositoryEnv.scientificFieldTagRepository,
	}
	renderService := &services.RenderService{
		BranchRepository:      repositoryEnv.branchRepository,
		PostRepository:        repositoryEnv.postRepository,
		ProjectPostRepository: repositoryEnv.projectPostRepository,
		Filesystem:            fs,
		BranchService:         nil, // Circular dependency filled in later...
	}
	postCollaboratorService := &services.PostCollaboratorService{
		PostCollaboratorRepository: repositoryEnv.postCollaboratorRepository,
		MemberRepository:           repositoryEnv.memberRepository,
		PostRepository:             repositoryEnv.postRepository,
	}
	branchCollaboratorService := &services.BranchCollaboratorService{
		BranchCollaboratorRepository: repositoryEnv.branchCollaboratorRepository,
		MemberRepository:             repositoryEnv.memberRepository,
	}
	postService := &services.PostService{
		PostRepository:                        repositoryEnv.postRepository,
		ProjectPostRepository:                 repositoryEnv.projectPostRepository,
		MemberRepository:                      repositoryEnv.memberRepository,
		ScientificFieldTagContainerRepository: repositoryEnv.scientificFieldTagContainerRepository,
		Filesystem:                            fs,
		PostCollaboratorService:               postCollaboratorService,
		RenderService:                         renderService,
		TagService:                            tagService,
	}
	branchService := &services.BranchService{
		BranchRepository:              repositoryEnv.branchRepository,
		ClosedBranchRepository:        repositoryEnv.closedBranchRepository,
		PostRepository:                repositoryEnv.postRepository,
		ProjectPostRepository:         repositoryEnv.projectPostRepository,
		ReviewRepository:              repositoryEnv.reviewRepository,
		DiscussionContainerRepository: repositoryEnv.discussionContainerRepository,
		DiscussionRepository:          repositoryEnv.discussionRepository,
		MemberRepository:              repositoryEnv.memberRepository,
		Filesystem:                    fs,
		RenderService:                 renderService,
		BranchCollaboratorService:     branchCollaboratorService,
		PostCollaboratorService:       postCollaboratorService,
		TagService:                    tagService,
	}
	projectPostService := &services.ProjectPostService{
		ProjectPostRepository:                 repositoryEnv.projectPostRepository,
		MemberRepository:                      repositoryEnv.memberRepository,
		ClosedBranchRepository:                repositoryEnv.closedBranchRepository,
		PostRepository:                        repositoryEnv.postRepository,
		ScientificFieldTagContainerRepository: repositoryEnv.scientificFieldTagContainerRepository,
		Filesystem:                            renderService.Filesystem,
		PostCollaboratorService:               postCollaboratorService,
		BranchCollaboratorService:             branchCollaboratorService,
		BranchService:                         branchService,
		TagService:                            tagService,
	}
	discussionService := &services.DiscussionService{
		DiscussionRepository:          repositoryEnv.discussionRepository,
		DiscussionContainerRepository: repositoryEnv.discussionContainerRepository,
		MemberRepository:              repositoryEnv.memberRepository,
	}
	renderService.BranchService = branchService // added afterwards since both require eachother

	// TODO we really need an automated DI solution..
	return ServiceEnv{
		postService: postService,
		memberService: &services.MemberService{
			MemberRepository: repositoryEnv.memberRepository,
		},
		branchService:             branchService,
		renderService:             renderService,
		projectPostService:        projectPostService,
		postCollaboratorService:   postCollaboratorService,
		branchCollaboratorService: branchCollaboratorService,
		discussionService:         discussionService,
		discussionContainerService: &services.DiscussionContainerService{
			DiscussionContainerRepository: repositoryEnv.discussionContainerRepository,
		},
		tagService: tagService,
		scientificFieldTagContainerService: &services.ScientificFieldTagContainerService{
			ContainerRepository: repositoryEnv.scientificFieldTagContainerRepository,
		},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController: &controllers.PostController{
			PostService:             serviceEnv.postService,
			RenderService:           serviceEnv.renderService,
			PostCollaboratorService: serviceEnv.postCollaboratorService,
		},
		memberController: &controllers.MemberController{
			MemberService: serviceEnv.memberService,
			TagService:    serviceEnv.tagService,
		},
		projectPostController: &controllers.ProjectPostController{
			ProjectPostService:         serviceEnv.projectPostService,
			DiscussionContainerService: serviceEnv.discussionContainerService,
			PostService:                serviceEnv.postService,
			RenderService:              serviceEnv.renderService,
		},
		discussionController: &controllers.DiscussionController{
			DiscussionService: serviceEnv.discussionService,
		},
		discussionContainerController: &controllers.DiscussionContainerController{
			DiscussionContainerService: serviceEnv.discussionContainerService,
		},
		filterController: &controllers.FilterController{
			PostService: serviceEnv.postService,
		},
		branchController: &controllers.BranchController{
			BranchService:             serviceEnv.branchService,
			RenderService:             serviceEnv.renderService,
			BranchCollaboratorService: serviceEnv.branchCollaboratorService,
		},
		tagController: &controllers.TagController{
			TagService:                         serviceEnv.tagService,
			ScientificFieldTagContainerService: serviceEnv.scientificFieldTagContainerService,
		},
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

	router := SetUpRouter(&controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
