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
	versionRepository database.RepositoryInterface[*models.Version]
	branchRepository  database.RepositoryInterface[*models.Branch]
	postRepository    database.RepositoryInterface[*models.Post]
}

type ServiceEnv struct {
	postService    interfaces.PostService
	versionService interfaces.VersionService
	memberService  interfaces.MemberService
	branchService  interfaces.BranchService
}

type ControllerEnv struct {
	postController        *controllers.PostController
	memberController      *controllers.MemberController
	projectPostController *controllers.ProjectPostController
	discussionController  *controllers.DiscussionController
	filterController      *controllers.FilterController
	branchController      *controllers.BranchController
	tagController         *controllers.TagController
	versionController     *controllers.VersionController
}

func initRepositoryEnv(db *gorm.DB) RepositoryEnv {
	return RepositoryEnv{
		versionRepository: &database.ModelRepository[*models.Version]{Database: db},
		branchRepository:  &database.ModelRepository[*models.Branch]{Database: db},
		postRepository:    &database.ModelRepository[*models.Post]{Database: db},
	}
}

func initServiceEnv(repositoryEnv RepositoryEnv, fs *filesystem.Filesystem) ServiceEnv {
	return ServiceEnv{
		postService: &services.PostService{},
		versionService: &services.VersionService{
			VersionRepository: repositoryEnv.versionRepository,
			Filesystem:        fs,
		},
		memberService: &services.MemberService{},
		branchService: &services.BranchService{
			PostRepository:   repositoryEnv.postRepository,
			BranchRepository: repositoryEnv.branchRepository,
			Filesystem:       fs,
		},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController:        &controllers.PostController{PostService: serviceEnv.postService},
		memberController:      &controllers.MemberController{MemberService: serviceEnv.memberService},
		projectPostController: &controllers.ProjectPostController{},
		discussionController:  &controllers.DiscussionController{},
		filterController:      &controllers.FilterController{},
		branchController:      &controllers.BranchController{BranchService: serviceEnv.branchService},
		tagController:         &controllers.TagController{},
		versionController:     &controllers.VersionController{VersionService: serviceEnv.versionService},
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
