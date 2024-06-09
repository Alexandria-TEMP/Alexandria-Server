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

	CollaborationType CollaborationType
}

type BranchCollaboratorDTO struct {
	ID                uint
	MemberID          uint
	BranchID          uint
	CollaborationType CollaborationType
}

func (model *BranchCollaborator) GetID() uint {
	return model.Model.ID
}

func (model *BranchCollaborator) IntoDTO() BranchCollaboratorDTO {
	return BranchCollaboratorDTO{
		model.ID,
		model.MemberID,
		model.BranchID,
		model.CollaborationType,
	}
}

func (model *BranchCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
