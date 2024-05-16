package models

import "gorm.io/gorm"

type ProjectPost struct {
	gorm.Model

	// ProjectPost belongs to Post
	Post   Post `gorm:"foreignKey:PostID"`
	PostID uint

	// ProjectPost has one ProjectMetadata
	ProjectMetadata ProjectMetadata `gorm:"foreignKey:ProjectPostID"`

	// ProjectPost has many MergeRequest
	OpenMergeRequests []MergeRequest `gorm:"foreignKey:ProjectPostID"`

	// ProjectPost has many ClosedMergeRequest
	ClosedMergeRequests []ClosedMergeRequest `gorm:"foreignKey:ProjectPostID"`
}

func (model *ProjectPost) GetID() uint {
	return model.Model.ID
}
