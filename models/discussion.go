package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Discussion struct {
	gorm.Model

	// Version has many Discussion
	VersionID uint

	// Discussion belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// Discussion optionally has many Discussion
	Replies  []Discussion `gorm:"foreignKey:ParentID"`
	ParentID *uint

	Text      string
	Deleted   bool
	Anonymous bool
}

func (model *Discussion) GetID() uint {
	return model.Model.ID
}

type DiscussionDTO struct {
	ID        uint
	VersionID uint
	MemberID  uint
	ParentID  *uint
	Replies   []DiscussionDTO
	Text      string
	Deleted   bool
	Anonymous bool
}

func (model *Discussion) MarshalJSON() ([]byte, error) {
	return json.Marshal(discussionIntoDTO(model))
}

// Helper function for MarshalJSON
func discussionIntoDTO(model *Discussion) DiscussionDTO {
	return DiscussionDTO{
		model.ID,
		model.VersionID,
		model.MemberID,
		model.ParentID,
		repliesIntoDTOs(model.Replies),
		model.Text,
		model.Deleted,
		model.Anonymous,
	}
}

// Helper function for MarshalJSON
func repliesIntoDTOs(replies []Discussion) []DiscussionDTO {
	result := make([]DiscussionDTO, len(replies))

	for i, reply := range replies {
		result[i] = discussionIntoDTO(&reply)
	}

	return result
}
