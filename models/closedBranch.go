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

	// SupercededBranch belongs to ClosedBranch
	SupercededBranch   Branch `gorm:"foreignKey:SupercededBranchID"`
	SupercededBranchID uint

	// ProjectPost has many ClosedBranch
	ProjectPostID uint

	BranchReviewDecision BranchReviewDecision
}

type ClosedBranchDTO struct {
	ID                   uint
	BranchID             uint
	SupercededBranchID   uint
	ProjectPostID        uint
	BranchReviewDecision BranchReviewDecision
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
