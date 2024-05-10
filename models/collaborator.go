package models

import "gorm.io/gorm"

type CollaborationType int16

const (
	Author CollaborationType = iota
	Contributor
	Reviewer
)

// A member that has collaborated on a post.
type PostCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member
	MemberID uint

	// PostMetadata has many PostCollaborator
	PostCollaboratorID uint

	CollaborationType CollaborationType
}

// A member that has collaborated on a merge request.
type MergeRequestCollaborator struct {
	gorm.Model

	// Belongs to Member
	Member   Member
	MemberID uint

	// MergeRequest has many MergeRequestCollaborator
	MergeRequestID uint

	CollaborationType CollaborationType
}
