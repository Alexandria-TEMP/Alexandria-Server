package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model

	// Post has one PostMetadata
	PostMetadata PostMetadata

	// Post belongs to Version
	CurrentVersion   Version
	CurrentVersionID uint
}
