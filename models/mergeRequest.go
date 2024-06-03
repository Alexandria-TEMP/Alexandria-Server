package models

import (
	"encoding/json"
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type MergeRequestReviewStatus string

const (
	MergeRequestOpenForReview MergeRequestReviewStatus = "open for review"
	MergeRequestPeerReviewed  MergeRequestReviewStatus = "peer reviewed"
	MergeRequestRejected      MergeRequestReviewStatus = "rejected"
)

type MergeRequest struct {
	gorm.Model

	/////////////////////////////////////////////
	// The MR's proposed changes:

	// MergeRequest belongs to Version
	NewVersion   Version `gorm:"foreignKey:NewVersionID"`
	NewVersionID uint

	UpdatedPostTitle string

	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields []*tags.ScientificFieldTag

	/////////////////////////////////////////////
	// The MR's metadata:

	// MergeRequest has many MergeRequestCollaborator
	Collaborators []*MergeRequestCollaborator `gorm:"foreignKey:MergeRequestID"`

	// MergeRequest has many MergeRequestReview
	Reviews []*MergeRequestReview `gorm:"foreignKey:MergeRequestID"`

	// ProjectPost has many MergeRequest
	ProjectPostID uint

	// MergeRequest belongs to Version (previous version)
	PreviousVersion   Version `gorm:"foreignKey:PreviousVersionID"`
	PreviousVersionID uint

	MergeRequestTitle    string
	Anonymous            bool
	MergeRequestDecision MergeRequestReviewStatus
}

type MergeRequestDTO struct {
	ID uint
	// MR's proposed changes
	NewVersionID                 uint
	NewPostTitle                 string
	UpdatedCompletionStatus      tags.CompletionStatus
	UpdatedScientificFieldTagIDs []uint
	// MR metadata
	CollaboratorIDs      []uint
	ReviewIDs            []uint
	ProjectPostID        uint
	PreviousVersionID    uint
	MergeRequestTitle    string
	Anonymous            bool
	MergeRequestDecision MergeRequestReviewStatus
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (model *MergeRequest) GetID() uint {
	return model.Model.ID
}

func (model *MergeRequest) IntoDTO() MergeRequestDTO {
	return MergeRequestDTO{
		model.ID,
		model.NewVersionID,
		model.UpdatedPostTitle,
		model.UpdatedCompletionStatus,
		tags.ScientificFieldTagIntoIDs(model.UpdatedScientificFields),
		mergeRequestCollaboratorsToIDs(model.Collaborators),
		reviewsToIDs(model.Reviews),
		model.ProjectPostID,
		model.PreviousVersionID,
		model.MergeRequestTitle,
		model.Anonymous,
		model.MergeRequestDecision,
		model.CreatedAt,
		model.UpdatedAt,
	}
}

func (model *MergeRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func mergeRequestCollaboratorsToIDs(collaborators []*MergeRequestCollaborator) []uint {
	ids := make([]uint, len(collaborators))

	for i, collaborator := range collaborators {
		ids[i] = collaborator.ID
	}

	return ids
}

// Helper function for JSON marshaling
func reviewsToIDs(reviews []*MergeRequestReview) []uint {
	ids := make([]uint, len(reviews))

	for i, review := range reviews {
		ids[i] = review.ID
	}

	return ids
}
