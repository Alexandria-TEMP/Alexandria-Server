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
	Replies  []*Discussion `gorm:"foreignKey:ParentID"`
	ParentID *uint

	Text      string
	Deleted   bool
	Anonymous bool
}

type DiscussionDTO struct {
	ID        uint
	VersionID uint
	MemberID  uint
	ReplyIDs  []uint
	Text      string
	Deleted   bool
	Anonymous bool
}

func (model *Discussion) GetID() uint {
	return model.Model.ID
}

func (model *Discussion) IntoDTO() DiscussionDTO {
	return DiscussionDTO{
		model.ID,
		model.VersionID,
		model.MemberID,
		discussionsIntoIDs(model.Replies),
		model.Text,
		model.Deleted,
		model.Anonymous,
	}
}

func (model *Discussion) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
