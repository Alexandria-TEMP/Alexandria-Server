package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type MergeRequestReviewDecision string

const (
	ReviewRejected MergeRequestReviewDecision = "rejected"
	ReviewApproved MergeRequestReviewDecision = "approved"
)

type MergeRequestReview struct {
	gorm.Model

	// MergeRequest has many MergeRequestReview
	MergeRequestID uint

	// MergeRequestReview belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	MergeRequestDecision MergeRequestReviewDecision
	Feedback             string
}

type MergeRequestReviewDTO struct {
	ID                   uint
	MergeRequestID       uint
	MemberID             uint
	MergeRequestDecision MergeRequestReviewDecision
	Feedback             string
	CreatedAt            time.Time
}

func (model *MergeRequestReview) GetID() uint {
	return model.Model.ID
}

func (model *MergeRequestReview) IntoDTO() MergeRequestReviewDTO {
	return MergeRequestReviewDTO{
		model.ID,
		model.MergeRequestID,
		model.MemberID,
		model.MergeRequestDecision,
		model.Feedback,
		model.CreatedAt,
	}
}

func (model *MergeRequestReview) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
