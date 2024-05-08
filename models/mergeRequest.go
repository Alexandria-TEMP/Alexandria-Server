package models

import (
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type MergeRequest struct {
	CreatedAt               time.Time
	UpdatedCompletionStatus tags.CompletionStatusTag
	UpdatedScientificFields tags.ScientificFieldTag
	NewVersion              Version
	Reviews                 []MergeRequestReview
	Collaborators           []Collaborator
	Anonymous               bool
}
