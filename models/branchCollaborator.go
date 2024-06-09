package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// A member that has collaborated on a branch.
type BranchCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// Branch has many BranchCollaborator
	BranchID uint
}

type BranchCollaboratorDTO struct {
	ID       uint `json:"id"`
	MemberID uint `json:"memberID"`
	BranchID uint `json:"branchID"`
}

func (model *BranchCollaborator) GetID() uint {
	return model.Model.ID
}

func (model *BranchCollaborator) IntoDTO() BranchCollaboratorDTO {
	return BranchCollaboratorDTO{
		model.ID,
		model.MemberID,
		model.BranchID,
	}
}

func (model *BranchCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
