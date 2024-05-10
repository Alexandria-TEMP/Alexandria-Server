package models

import (
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type MergeRequest struct {
	gorm.Model

	// MergeRequest belongs to Version
	NewVersion   Version
	NewVersionID uint

	// MergeRequest has many MergeRequestCollaborator
	Collaborators []MergeRequestCollaborator

	// MergeRequest has many MergeRequestReview
	Reviews []MergeRequestReview

	// ProjectPost has many MergeRequest
	ProjectPostID uint

	CreatedAt               time.Time
	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields tags.ScientificField
	Anonymous               bool
}
