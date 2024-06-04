package models

import (
	"encoding/json"

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

	MergeRequestDecision MergeRequestReviewDecision
}

type ClosedMergeRequestDTO struct {
	ID                      uint
	MergeRequestID          uint
	MainVersionWhenClosedID uint
	ProjectPostID           uint
	MergeRequestDecision    MergeRequestReviewDecision
}

func (model *ClosedMergeRequest) GetID() uint {
	return model.Model.ID
}

func (model *ClosedMergeRequest) IntoDTO() ClosedMergeRequestDTO {
	return ClosedMergeRequestDTO{
		model.ID,
		model.MergeRequestID,
		model.MainVersionWhenClosedID,
		model.ProjectPostID,
		model.MergeRequestDecision,
	}
}

func (model *ClosedMergeRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
