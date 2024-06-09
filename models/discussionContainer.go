package models

import "gorm.io/gorm"

type DiscussionContainer struct {
	gorm.Model

	// DiscussionContainer has many Discussion
	Discussions []*Discussion `gorm:"foreignKey:ContainerID"`
}

func (model *DiscussionContainer) GetID() uint {
	return model.Model.ID
}
