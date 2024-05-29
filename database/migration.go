package database

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gorm.io/gorm"
)

// Changes to the models require database migrations. All models should be migrated here.
func autoMigrateAllModels(db *gorm.DB) error {
	// NOTE FOR FUTURE CHANGES: the order of migrations matters!
	// Foreign keys (e.g. uint) have to be initialized AFTER the
	// model that is being pointed to has been migrated.
	// For example, if "Foo has one Bar" (meaning Foo holds "Bar",
	// and Bar holds "FooID uint"), Foo should be migrated before Bar.
	//
	// If this is not upheld, foreign key constraint errors will be thrown.
	return db.AutoMigrate(
		&models.Version{},            //
		&models.Post{},               // FK to Version
		&models.ProjectPost{},        // FK to Post
		&models.Branch{},             // FK to Version, ProjectPost
		&models.ClosedBranch{},       // FK to Branch, Version, ProjectPost
		&models.Member{},             //
		&models.PostCollaborator{},   // FK to Member, PostMetadata
		&models.BranchCollaborator{}, // FK to Member, Branch
		&models.Discussion{},         // FK to Version, Member
		&models.BranchReview{},       // FK to Branch, Member
	)
}
