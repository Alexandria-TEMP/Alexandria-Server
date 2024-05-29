package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// A branch that is no longer open, including merged and non-merged.
type ClosedBranch struct {
	gorm.Model

	// ClosedBranch belongs to Branch
	Branch   Branch `gorm:"foreignKey:BranchID"`
	BranchID uint

	// ClosedBranch belongs to Version
	MainVersionWhenClosed   Version `gorm:"foreignKey:MainVersionWhenClosedID"`
	MainVersionWhenClosedID uint

	// ProjectPost has many ClosedBranch
	ProjectPostID uint

	BranchDecision BranchDecision
}

type ClosedBranchDTO struct {
	ID                      uint
	BranchID                uint
	MainVersionWhenClosedID uint
	ProjectPostID           uint
	BranchDecision          BranchDecision
}

func (model *ClosedBranch) GetID() uint {
	return model.Model.ID
}

func (model *ClosedBranch) IntoDTO() ClosedBranchDTO {
	return ClosedBranchDTO{
		model.ID,
		model.BranchID,
		model.MainVersionWhenClosedID,
		model.ProjectPostID,
		model.BranchDecision,
	}
}

func (model *ClosedBranch) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
