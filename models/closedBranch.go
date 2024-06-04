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

	// SupercededBranch belongs to ClsoedBranch
	SupercededBranch   Branch `gorm:"foreignKey:SupercededBranchID"`
	SupercededBranchID uint

	// ClosedBranch belongs to ProjectPost
	ProjectPost   ProjectPost `gorm:"foreignKey:ProjectPostID"`
	ProjectPostID uint

	BranchDecision BranchDecision
}

type ClosedBranchDTO struct {
	ID                 uint
	BranchID           uint
	SupercededBranchID uint
	ProjectPostID      uint
	BranchDecision     BranchDecision
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
		model.BranchDecision,
	}
}

func (model *ClosedBranch) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
