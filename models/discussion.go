package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Discussion struct {
	gorm.Model

	// Discussion belongs to DiscussionContainer
	ContainerID uint

	// Discussion optionally belongs to Member (anonymity is possible)
	Member   *Member `gorm:"foreignKey:MemberID"`
	MemberID *uint

	// Discussion optionally has many Discussion
	Replies  []*Discussion `gorm:"foreignKey:ParentID"`
	ParentID *uint

	Text string
}

type DiscussionDTO struct {
	ID       uint   `json:"id" example:"1"`
	MemberID *uint  `json:"memberID" example:"1"`
	ReplyIDs []uint `json:"replyIDs" example:"2"`
	Text     string `json:"text" example:"Discussion content."`
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
