package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// A branch that is no longer open, including merged and non-merged.
type ClosedBranch struct {
	gorm.Model

	// Branch belongs to ClosedBranch
	Branch   Branch `gorm:"foreignKey:BranchID"`
	BranchID uint

	// SupercededBranch may belong to ClosedBranch
	SupercededBranch   *Branch `gorm:"foreignKey:SupercededBranchID"`
	SupercededBranchID *uint

	// ProjectPost has many ClosedBranch
	ProjectPostID uint

	BranchReviewDecision BranchReviewDecision
}

type ClosedBranchDTO struct {
	ID                   uint                 `json:"id" example:"1"`
	BranchID             uint                 `json:"branchID" example:"1"`
	SupercededBranchID   *uint                `json:"supercededBranchID" example:"2"`
	ProjectPostID        uint                 `json:"projectPostID" example:"1"`
	BranchReviewDecision BranchReviewDecision `json:"branchReviewDecision" example:"approved"`
}

func (model *ClosedBranch) GetID() uint {
	return model.Model.ID
}

func (model *ClosedBranch) IntoDTO() ClosedBranchDTO {
	return ClosedBranchDTO{
		model.ID,
		model.BranchID,
		model.SupercededBranchID,
		model.ProjectPostID,
		model.BranchReviewDecision,
	}
}

func (model *ClosedBranch) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
