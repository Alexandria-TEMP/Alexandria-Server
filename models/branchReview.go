package models

import (
	"encoding/json"
	"slices"
	"time"

	"gorm.io/gorm"
)

type BranchReviewDecision string

const (
	Rejected BranchReviewDecision = "rejected"
	Approved BranchReviewDecision = "approved"
)

func (enum *BranchReviewDecision) IsValid() bool {
	valid := []BranchReviewDecision{Rejected, Approved}
	return slices.Contains(valid, *enum)
}

type BranchReview struct {
	gorm.Model

	// Branch has many BranchReview
	BranchID uint

	// BranchReview belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	BranchReviewDecision BranchReviewDecision
	Feedback             string
}

type BranchReviewDTO struct {
	ID                   uint                 `json:"id" example:"1"`
	BranchID             uint                 `json:"branchID" example:"1"`
	MemberID             uint                 `json:"memberID" example:"1"`
	BranchReviewDecision BranchReviewDecision `json:"branchReviewDecision" example:"approved"`
	Feedback             string               `json:"feedback" example:"Fantastic work!"`
	CreatedAt            time.Time            `json:"createdAt"`
}

func (model *BranchReview) GetID() uint {
	return model.Model.ID
}

func (model *BranchReview) IntoDTO() BranchReviewDTO {
	return BranchReviewDTO{
		model.ID,
		model.BranchID,
		model.MemberID,
		model.BranchReviewDecision,
		model.Feedback,
		model.CreatedAt,
	}
}

func (model *BranchReview) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
