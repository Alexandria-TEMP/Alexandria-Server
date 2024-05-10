package models

import (
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type MergeRequest struct {
	gorm.Model
	CreatedAt               time.Time
	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields tags.ScientificField
	NewVersion              Version
	Reviews                 []MergeRequestReview
	Collaborators           []Collaborator
	Anonymous               bool
}
