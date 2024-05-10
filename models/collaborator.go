package models

import "gorm.io/gorm"

type CollaborationType string

const (
	Author      CollaborationType = "author"
	Contributor CollaborationType = "contributor"
	Reviewer    CollaborationType = "reviewer"
)

// A member that has collaborated on a post.
type PostCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// PostMetadata has many PostCollaborator
	PostMetadataID uint

	CollaborationType CollaborationType `gorm:"serializer:json"`
}

// A member that has collaborated on a merge request.
type MergeRequestCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// MergeRequest has many MergeRequestCollaborator
	MergeRequestID uint

	CollaborationType CollaborationType `gorm:"serializer:json"`
}
