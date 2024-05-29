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

	// ProjectPost has many Branch
	OpenBranchs []*Branch `gorm:"foreignKey:ProjectPostID"`

	// ProjectPost has many ClosedBranch
	ClosedBranchs []*ClosedBranch `gorm:"foreignKey:ProjectPostID"`

	CompletionStatus    tags.CompletionStatus
	FeedbackPreference  tags.FeedbackPreference
	PostReviewStatusTag tags.PostReviewStatus
}

type ProjectPostDTO struct {
	ID                  uint
	PostDTO             PostDTO
	OpenBranchIDs       []uint
	ClosedBranchIDs     []uint
	CompletionStatus    tags.CompletionStatus
	FeedbackPreference  tags.FeedbackPreference
	PostReviewStatusTag tags.PostReviewStatus
}

func (model *ProjectPost) GetID() uint {
	return model.Model.ID
}

func (model *ProjectPost) IntoDTO() ProjectPostDTO {
	return ProjectPostDTO{
		model.ID,
		model.Post.IntoDTO(),
		branchsToIDs(model.OpenBranchs),
		closedBranchsToIDs(model.ClosedBranchs),
		model.CompletionStatus,
		model.FeedbackPreference,
		model.PostReviewStatusTag,
	}
}

func (model *ProjectPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func branchsToIDs(branchs []*Branch) []uint {
	ids := make([]uint, len(branchs))

	for i, branchs := range branchs {
		ids[i] = branchs.ID
	}

	return ids
}

// Helper function for JSON marshaling
func closedBranchsToIDs(branchs []*ClosedBranch) []uint {
	ids := make([]uint, len(branchs))

	for i, branchs := range branchs {
		ids[i] = branchs.ID
	}

	return ids
}
