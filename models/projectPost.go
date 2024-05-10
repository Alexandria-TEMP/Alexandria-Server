package models

import "gorm.io/gorm"

type ProjectPost struct {
	gorm.Model

	// ProjectPost belongs to Post
	Post   Post
	PostID uint

	// ProjectPost has one ProjectMetadata
	ProjectMetadata ProjectMetadata

	// ProjectPost has many MergeRequest
	OpenMergeRequests []MergeRequest

	// ProjectPost has many ClosedMergeRequest
	ClosedMergeRequests []ClosedMergeRequest
}
