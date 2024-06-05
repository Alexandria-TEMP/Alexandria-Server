package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type BranchDecision string

const (
	Rejected BranchDecision = "rejected"
	Approved BranchDecision = "approved"
)

type Review struct {
	gorm.Model

	// Branch has many Review
	BranchID uint

	// Review belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	BranchDecision BranchDecision
	Feedback       string
}

type ReviewDTO struct {
	ID             uint
	BranchID       uint
	MemberID       uint
	BranchDecision BranchDecision
	Feedback       string
	CreatedAt      time.Time
}

func (model *Review) GetID() uint {
	return model.Model.ID
}

func (model *Review) IntoDTO() ReviewDTO {
	return ReviewDTO{
		model.ID,
		model.BranchID,
		model.MemberID,
		model.BranchDecision,
		model.Feedback,
		model.CreatedAt,
	}
}

func (model *Review) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
