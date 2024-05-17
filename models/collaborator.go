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

func (model *PostCollaborator) GetID() uint {
	return model.Model.ID
}

type PostCollaboratorDTO struct {
	ID                uint
	MemberID          uint
	PostID            uint
	CollaborationType CollaborationType
}

func (model *PostCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(PostCollaboratorDTO{
		model.ID,
		model.MemberID,
		model.PostID,
		model.CollaborationType,
	})
}

// A member that has collaborated on a merge request.
type MergeRequestCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// MergeRequest has many MergeRequestCollaborator
	MergeRequestID uint

	CollaborationType CollaborationType
}

func (model *MergeRequestCollaborator) GetID() uint {
	return model.Model.ID
}

type MergeRequestCollaboratorDTO struct {
	ID                uint
	MemberID          uint
	PostID            uint
	CollaborationType CollaborationType
}

func (model *MergeRequestCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeRequestCollaboratorDTO{
		model.ID,
		model.MemberID,
		model.MergeRequestID,
		model.CollaborationType,
	})
}
