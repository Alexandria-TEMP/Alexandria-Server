package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type RenderStatus string

const (
	RenderSuccess RenderStatus = "success"
	RenderPending RenderStatus = "pending"
	RenderFailure RenderStatus = "failure"
)

type Version struct {
	gorm.Model

	// Version has many Discussion
	Discussions  []*Discussion `gorm:"foreignKey:VersionID"`
	RenderStatus RenderStatus
}

type VersionDTO struct {
	ID            uint
	DiscussionIDs []uint
	RenderStatus  RenderStatus
}

func (model *Version) GetID() uint {
	return model.Model.ID
}

func (model *Version) IntoDTO() VersionDTO {
	return VersionDTO{
		model.ID,
		discussionsIntoIDs(model.Discussions),
		model.RenderStatus,
	}
}

func (model *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func discussionsIntoIDs(discussions []*Discussion) []uint {
	ids := make([]uint, len(discussions))

	for i, discussion := range discussions {
		ids[i] = discussion.ID
	}

	return ids
}
