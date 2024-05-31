package server

import (
	"log"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gorm.io/gorm"
)

type RepositoryEnv struct {
	memberRepository 	database.RepositoryInterface[*models.Member]
	versionRepository 	database.RepositoryInterface[*models.Version]
	tagRepository		database.RepositoryInterface[*tags.ScientificFieldTag]
}

type ServiceEnv struct {
	postService    services.PostService
	versionService services.VersionService
	memberService  services.MemberService
	tagService     services.TagService
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

func initRepositoryEnv(db *gorm.DB) RepositoryEnv {
	return RepositoryEnv{
		memberRepository: &database.ModelRepository[*models.Member]{Database: db},
		//memberRepository: database.ModelRepository[*models.Member]{Database: db},
		versionRepository: &database.ModelRepository[*models.Version]{Database: db},
		tagRepository: &database.ModelRepository[*tags.ScientificFieldTag]{Database: db},
	}
}


func initServiceEnv(repositoryEnv *RepositoryEnv, fs *filesystem.Filesystem) ServiceEnv {
	return ServiceEnv{
		postService:    services.PostService{},
		versionService: services.VersionService{
			VersionRepository: repositoryEnv.versionRepository,
			Filesystem: fs},
		memberService:  services.MemberService{
			MemberRepository: repositoryEnv.memberRepository,
		},
	}
}

func initControllerEnv(serviceEnv *ServiceEnv) ControllerEnv {
	return ControllerEnv{
		postController:         controllers.PostController{PostService: &serviceEnv.postService},
		memberController:       controllers.MemberController {
			MemberService: &serviceEnv.memberService, 
			TagService: &serviceEnv.tagService,
		},
		projectPostController:  controllers.ProjectPostController{},
		discussionController:   controllers.DiscussionController{},
		filterController:       controllers.FilterController{},
		mergeRequestController: controllers.MergeRequestController{},
		tagController:          controllers.TagController{},
		versionController:      controllers.VersionController{},
	}
}

func Init() {
	db, err := database.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.InitFilesystem()

	repositoryEnv := initRepositoryEnv(db)
	serviceEnv := initServiceEnv(&repositoryEnv, fs)
	controllerEnv := initControllerEnv(&serviceEnv)

	router := SetUpRouter(&controllerEnv)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
