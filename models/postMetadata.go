package models

import (
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type PostMetadata struct {
	gorm.Model

	// PostMetadata has many PostCollaborator
	Collaborators []PostCollaborator `gorm:"foreignKey:PostMetadataID"`

	// Post has one PostMetadata
	PostID uint

	CreatedAt           time.Time
	UpdatedAt           time.Time
	PostType            tags.PostType
	ScientificFieldTags []tags.ScientificField `gorm:"serializer:json"`
}

func (model *PostMetadata) GetID() uint {
	return model.Model.ID
}
