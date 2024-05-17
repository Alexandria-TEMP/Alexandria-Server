package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository struct {
	// TODO write serialization/deserialization, OR use a filesystem instead
	// QuartoProject multipart.File `swaggerignore:"true"`
}

type Version struct {
	gorm.Model

	Repository Repository `gorm:"serializer:json"`

	// Version has many Discussion
	Discussions []*Discussion `gorm:"foreignKey:VersionID"`
}

func (model *Version) GetID() uint {
	return model.Model.ID
}

type VersionDTO struct {
	ID            uint
	DiscussionIDs []uint
}

func (model *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(VersionDTO{
		model.ID,
		discussionsIntoIDs(model.Discussions),
	})
}

// Helper function for JSON marshaling
func discussionsIntoIDs(discussions []*Discussion) []uint {
	ids := make([]uint, len(discussions))

	for i, discussion := range discussions {
		ids[i] = discussion.ID
	}

	return ids
}
