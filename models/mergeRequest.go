package models

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type MergeRequest struct {
	gorm.Model

	// MergeRequest belongs to Version
	NewVersion   Version `gorm:"foreignKey:NewVersionID"`
	NewVersionID uint

	// MergeRequest has many MergeRequestCollaborator
	Collaborators []*MergeRequestCollaborator `gorm:"foreignKey:MergeRequestID"`

	// MergeRequest has many MergeRequestReview
	Reviews []*MergeRequestReview `gorm:"foreignKey:MergeRequestID"`

	// ProjectPost has many MergeRequest
	ProjectPostID uint

	Title                   string
	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields tags.ScientificField `gorm:"serializer:json"`
	Anonymous               bool
}

type MergeRequestDTO struct {
	ID                      uint
	NewVersionID            uint
	CollaboratorIDs         []uint
	ReviewIDs               []uint
	ProjectPostID           uint
	Title                   string
	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields tags.ScientificField
	Anonymous               bool
}

func (model *MergeRequest) GetID() uint {
	return model.Model.ID
}

func (model *MergeRequest) IntoDTO() MergeRequestDTO {
	return MergeRequestDTO{
		model.ID,
		model.NewVersionID,
		mergeRequestCollaboratorsToIDs(model.Collaborators),
		reviewsToIDs(model.Reviews),
		model.ProjectPostID,
		model.Title,
		model.UpdatedCompletionStatus,
		model.UpdatedScientificFields,
		model.Anonymous,
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
