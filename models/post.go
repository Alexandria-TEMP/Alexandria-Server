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
	ID               uint
	CurrentVersionID uint
	// TODO add fields
}

func (model *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(PostDTO{
		// TODO add fields
	})
}
