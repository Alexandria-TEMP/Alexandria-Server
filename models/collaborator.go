package models

import (
	"encoding/json"
	"slices"

	"gorm.io/gorm"
)

type CollaborationType string

const (
	Author      CollaborationType = "author"
	Contributor CollaborationType = "contributor"
	Reviewer    CollaborationType = "reviewer"
)

func (enum *CollaborationType) IsValid() bool {
	valid := []CollaborationType{Author, Contributor, Reviewer}
	return slices.Contains(valid, *enum)
}

// A member that has collaborated on a post.
type PostCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// Post has many PostCollaborator
	PostID uint

	CollaborationType CollaborationType
}

type PostCollaboratorDTO struct {
	ID                uint
	MemberID          uint
	PostID            uint
	CollaborationType CollaborationType
}

func (model *PostCollaborator) GetID() uint {
	return model.Model.ID
}

func (model *PostCollaborator) IntoDTO() PostCollaboratorDTO {
	return PostCollaboratorDTO{
		model.ID,
		model.MemberID,
		model.PostID,
		model.CollaborationType,
	}
}

func (model *PostCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

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
	ID       uint
	MemberID uint
	BranchID uint
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
