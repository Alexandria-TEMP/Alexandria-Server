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
	OpenBranches []*Branch `gorm:"foreignKey:ProjectPostID"`

	// ProjectPost has many ClosedBranch
	ClosedBranches []*ClosedBranch `gorm:"foreignKey:ProjectPostID"`

	CompletionStatus    tags.CompletionStatus
	FeedbackPreference  tags.FeedbackPreference
	PostReviewStatusTag tags.PostReviewStatus

	// ProjectPost may have a Branch (it won't if its master)
	LastMergedBranch *Branch `gorm:"foreignKey:ProjectPostID"`
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
		branchesToIDs(model.OpenBranches),
		closedBranchesToIDs(model.ClosedBranches),
		model.CompletionStatus,
		model.FeedbackPreference,
		model.PostReviewStatusTag,
	}
}

func (model *ProjectPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func branchesToIDs(branches []*Branch) []uint {
	ids := make([]uint, len(branches))

	for i, branches := range branches {
		ids[i] = branches.ID
	}

	return ids
}

// Helper function for JSON marshaling
func closedBranchesToIDs(branches []*ClosedBranch) []uint {
	ids := make([]uint, len(branches))

	for i, branches := range branches {
		ids[i] = branches.ID
	}

	return ids
}
