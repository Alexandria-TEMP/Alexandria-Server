package models

import "gorm.io/gorm"

type ProjectPost struct {
	gorm.Model
	Post
	ProjectMetadata
	OpenMergeRequests   []MergeRequest
	ClosedMergeRequests []ClosedMergeRequest
}
