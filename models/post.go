package models

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	// Post has many PostCollaborator
	Collaborators []PostCollaborator `gorm:"foreignKey:PostID"`

	// Post belongs to Version
	CurrentVersion   Version `gorm:"foreignKey:CurrentVersionID"`
	CurrentVersionID uint

	PostType            tags.PostType
	ScientificFieldTags []tags.ScientificField `gorm:"serializer:json"`
}

func (model *Post) GetID() uint {
	return model.Model.ID
}

type PostDTO struct {
	ID                  uint
	CollaboratorIDs     []uint
	VersionID           uint
	PostType            tags.PostType
	ScientificFieldTags []tags.ScientificField
}

func (model *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(PostDTO{
		model.ID,
		postCollaboratorsToIDs(model.Collaborators),
		model.CurrentVersionID,
		model.PostType,
		model.ScientificFieldTags,
	})
}

// Helper function for JSON marshaling
func postCollaboratorsToIDs(collaborators []PostCollaborator) []uint {
	ids := make([]uint, len(collaborators))

	for i, collaborator := range collaborators {
		ids[i] = collaborator.ID
	}

	return ids
}
