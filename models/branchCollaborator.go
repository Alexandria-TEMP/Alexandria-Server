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
	ID       uint `json:"id" example:"1"`
	MemberID uint `json:"memberID" example:"1"`
	BranchID uint `json:"branchID" example:"1"`
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
