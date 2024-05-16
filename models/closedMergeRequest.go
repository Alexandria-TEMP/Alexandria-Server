package models

import (
	"time"

	"gorm.io/gorm"
)

// A merge request that is no longer open, including merged and non-merged.
type ClosedMergeRequest struct {
	gorm.Model

	// ClosedMergeRequest belongs to MergeRequest
	MergeRequest   MergeRequest `gorm:"foreignKey:MergeRequestID"`
	MergeRequestID uint

	// ClosedMergeRequest belongs to Version
	MainVersionWhenClosed   Version `gorm:"foreignKey:MainVersionWhenClosedID"`
	MainVersionWhenClosedID uint

	// ProjectPost has many ClosedMergeRequest
	ProjectPostID uint

	CreatedAt            time.Time
	MergeRequestDecision MergeRequestDecision
}

func (model *ClosedMergeRequest) GetID() uint {
	return model.Model.ID
}
