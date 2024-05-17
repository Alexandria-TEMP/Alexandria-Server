package models

import (
	"encoding/json"

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

type ProjectPostDTO struct {
	ID                    uint
	PostID                uint
	OpenMergeRequestIDs   []uint
	ClosedMergeRequestIDs []uint
	CompletionStatus      tags.CompletionStatus
	FeedbackPreference    tags.FeedbackPreference
	PostReviewStatusTag   tags.PostReviewStatus
}

func (model *ProjectPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(ProjectPostDTO{
		model.ID,
		model.PostID,
		mergeRequestsToIDs(model.OpenMergeRequests),
		closedMergeRequestsToIDs(model.ClosedMergeRequests),
		model.CompletionStatus,
		model.FeedbackPreference,
		model.PostReviewStatusTag,
	})
}

// Helper function for JSON marshaling
func mergeRequestsToIDs(mergeRequests []MergeRequest) []uint {
	ids := make([]uint, len(mergeRequests))

	for i, mergeRequests := range mergeRequests {
		ids[i] = mergeRequests.ID
	}

	return ids
}

// Helper function for JSON marshaling
func closedMergeRequestsToIDs(mergeRequests []ClosedMergeRequest) []uint {
	ids := make([]uint, len(mergeRequests))

	for i, mergeRequests := range mergeRequests {
		ids[i] = mergeRequests.ID
	}

	return ids
}
