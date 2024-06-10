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

type DiscussionContainerDTO struct {
	ID            uint   `json:"id"`
	DiscussionIDs []uint `json:"discussionIDs"`
}

func (model *DiscussionContainer) GetID() uint {
	return model.Model.ID
}

func (model *DiscussionContainer) IntoDTO() DiscussionContainerDTO {
	return DiscussionContainerDTO{
		ID:            model.ID,
		DiscussionIDs: discussionsIntoIDs(model.Discussions),
	}
}

func (model *DiscussionContainer) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
