package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type CollaborationType string

const (
	Author      CollaborationType = "author"
	Contributor CollaborationType = "contributor"
	Reviewer    CollaborationType = "reviewer"
)

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

// A member that has collaborated on a merge request.
type MergeRequestCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// MergeRequest has many MergeRequestCollaborator
	MergeRequestID uint

	// Merge request collaborators don't have a collaboration type,
	// because there is no concept of contributor/reviewer.
}

type MergeRequestCollaboratorDTO struct {
	ID             uint
	MemberID       uint
	MergeRequestID uint
}

func (model *MergeRequestCollaborator) GetID() uint {
	return model.Model.ID
}

func (model *MergeRequestCollaborator) IntoDTO() MergeRequestCollaboratorDTO {
	return MergeRequestCollaboratorDTO{
		model.ID,
		model.MemberID,
		model.MergeRequestID,
	}
}

func (model *MergeRequestCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
