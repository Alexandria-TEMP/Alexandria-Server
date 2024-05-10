package database

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gorm.io/gorm"
)

// Changes to the models require database migrations. All models should be migrated here.
func AutoMigrateAllModels(db *gorm.DB) error {
	// Listed in alphabetical order
	return db.AutoMigrate(
		&models.ClosedMergeRequest{},
		&models.Collaborator{},
		&models.Discussion{},
		&models.Member{},
		&models.MergeRequest{},
		&models.MergeRequestReview{},
		&models.Post{},
		&models.PostMetadata{},
		&models.ProjectMetadata{},
		&models.ProjectPost{},
		&models.Repository{},
		&models.Version{},
	)
}
