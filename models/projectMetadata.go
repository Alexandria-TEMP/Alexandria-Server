package models

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type ProjectMetadata struct {
	gorm.Model

	// ProjectPost has one ProjectMetadata
	ProjectPostID uint

	// TODO why is ForkedFrom a ClosedMergeRequest?
	// ForkedFrom          ClosedMergeRequest

	CompletionStatus    tags.CompletionStatus   `gorm:"serializer:json"`
	FeedbackPreference  tags.FeedbackPreference `gorm:"serializer:json"`
	PostReviewStatusTag tags.PostReviewStatus   `gorm:"serializer:json"`
}
