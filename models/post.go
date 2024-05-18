package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	// Post has one PostMetadata
	PostMetadata PostMetadata `gorm:"foreignKey:PostID"`

	// Post belongs to Version
	CurrentVersion   Version `gorm:"foreignKey:CurrentVersionID"`
	CurrentVersionID uint
}

func (model *Post) GetID() uint {
	return model.Model.ID
}
