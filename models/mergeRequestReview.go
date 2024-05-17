package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type MergeRequestDecision string

const (
	Rejected MergeRequestDecision = "rejected"
	Approved MergeRequestDecision = "approved"
)

type MergeRequestReview struct {
	gorm.Model

	// MergeRequest has many MergeRequestReview
	MergeRequestID uint

	// MergeRequestReview belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	MergeRequestDecision MergeRequestDecision
	Feedback             string
}

func (model *MergeRequestReview) GetID() uint {
	return model.Model.ID
}

type MergeRequestReviewDTO struct {
	ID                   uint
	MergeRequestID       uint
	MemberID             uint
	MergeRequestDecision MergeRequestDecision
	Feedback             string
}

func (model *MergeRequestReview) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeRequestReviewDTO{
		model.ID,
		model.MergeRequestID,
		model.MemberID,
		model.MergeRequestDecision,
		model.Feedback,
	})
}
