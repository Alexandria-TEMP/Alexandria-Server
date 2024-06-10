package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type DiscussionContainer struct {
	gorm.Model

	// DiscussionContainer has many Discussion
	Discussions []*Discussion `gorm:"foreignKey:ContainerID"`
}

func (model *DiscussionContainer) GetID() uint {
	return model.Model.ID
}

type Discussion struct {
	gorm.Model

	// Discussion belongs to DiscussionContainer
	ContainerID uint

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
	ID        uint   `json:"id"`
	MemberID  uint   `json:"memberID"`
	ReplyIDs  []uint `json:"replyIDs"`
	Text      string `json:"text"`
	Deleted   bool   `json:"deleted"`
	Anonymous bool   `json:"anonymous"`
}

func (model *Discussion) GetID() uint {
	return model.Model.ID
}

func (model *Discussion) IntoDTO() DiscussionDTO {
	return DiscussionDTO{
		model.ID,
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

// Helper function for JSON marshaling
func discussionsIntoIDs(discussions []*Discussion) []uint {
	ids := make([]uint, len(discussions))

	for i, discussion := range discussions {
		ids[i] = discussion.ID
	}

	return ids
}
