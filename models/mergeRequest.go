package models

import (
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type MergeRequest struct {
	NewVersion    Version
	Reviews       []MergeRequestReview
	Anonymous     bool
	CreatedAt     time.Time
	Collaborators []Collaborator
	// If this is nil (the defualt value) it means that there is no update
	UpdatedCompletionStatus tags.CompletionStatusTag
	// If this is nil (the defualt value) it means that there is no update
	UpdatedScientificFields tags.ScientificFieldTag
}
