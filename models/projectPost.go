package models

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type ProjectPost struct {
	gorm.Model

	// ProjectPost belongs to Post
	Post   Post `gorm:"foreignKey:PostID"`
	PostID uint

	// ProjectPost has many MergeRequest
	OpenMergeRequests []MergeRequest `gorm:"foreignKey:ProjectPostID"`

	// ProjectPost has many ClosedMergeRequest
	ClosedMergeRequests []ClosedMergeRequest `gorm:"foreignKey:ProjectPostID"`

	CompletionStatus    tags.CompletionStatus
	FeedbackPreference  tags.FeedbackPreference
	PostReviewStatusTag tags.PostReviewStatus
}

func (model *ProjectPost) GetID() uint {
	return model.Model.ID
}
