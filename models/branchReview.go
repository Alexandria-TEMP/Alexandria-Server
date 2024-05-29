package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type BranchDecision string

const (
	Rejected BranchDecision = "rejected"
	Approved BranchDecision = "approved"
)

type BranchReview struct {
	gorm.Model

	// Branch has many BranchReview
	BranchID uint

	// BranchReview belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	BranchDecision BranchDecision
	Feedback       string
}

type BranchReviewDTO struct {
	ID             uint
	BranchID       uint
	MemberID       uint
	BranchDecision BranchDecision
	Feedback       string
}

func (model *BranchReview) GetID() uint {
	return model.Model.ID
}

func (model *BranchReview) IntoDTO() BranchReviewDTO {
	return BranchReviewDTO{
		model.ID,
		model.BranchID,
		model.MemberID,
		model.BranchDecision,
		model.Feedback,
	}
}

func (model *BranchReview) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
