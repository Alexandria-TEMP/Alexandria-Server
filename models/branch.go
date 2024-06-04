package models

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type RenderStatus string

const (
	Success RenderStatus = "success"
	Pending RenderStatus = "pending"
	Failure RenderStatus = "failure"
)

type BranchReviewStatus string

const (
	BranchOpenForReview BranchReviewStatus = "open for review"
	BranchPeerReviewed  BranchReviewStatus = "peer reviewed"
	BranchRejected      BranchReviewStatus = "rejected"
)

type Branch struct {
	gorm.Model

	/////////////////////////////////////////////
	// The branch's proposed changes:

	NewPostTitle string

	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields []tags.ScientificField `gorm:"serializer:json"`

	/////////////////////////////////////////////
	// The branch's metadata:

	// Branch has many BranchCollaborator
	Collaborators []*BranchCollaborator `gorm:"foreignKey:BranchID"`

	// Branch has many BranchReview
	Reviews []*BranchReview `gorm:"foreignKey:BranchID"`

	// Branch has many Discussion
	Discussions []*Discussion `gorm:"foreignKey:BranchID"`

	// Branch has a ProjectPost
	ProjectPost   ProjectPost `gorm:"foreignKey:ProjectPostID"`
	ProjectPostID uint

	BranchTitle string

	Anonymous bool

	RenderStatus       RenderStatus
	BranchReviewStatus BranchReviewStatus
}

type BranchDTO struct {
	ID uint
	// MR's proposed changes
	NewPostTitle            string
	UpdatedCompletionStatus tags.CompletionStatus
	UpdatedScientificFields []tags.ScientificField
	// MR metadata
	CollaboratorIDs    []uint
	ReviewIDs          []uint
	ProjectPostID      uint
	BranchTitle        string
	Anonymous          bool
	RenderStatus       RenderStatus
	DiscussionIDs      []uint
	BranchReviewStatus BranchReviewStatus
}

func (model *Branch) GetID() uint {
	return model.Model.ID
}

func (model *Branch) IntoDTO() BranchDTO {
	return BranchDTO{
		model.ID,
		model.NewPostTitle,
		model.UpdatedCompletionStatus,
		model.UpdatedScientificFields,
		branchCollaboratorsToIDs(model.Collaborators),
		reviewsToIDs(model.Reviews),
		model.ProjectPostID,
		model.BranchTitle,
		model.Anonymous,
		model.RenderStatus,
		discussionsIntoIDs(model.Discussions),
		model.BranchReviewStatus,
	}
}

func (model *Branch) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func branchCollaboratorsToIDs(collaborators []*BranchCollaborator) []uint {
	ids := make([]uint, len(collaborators))

	for i, collaborator := range collaborators {
		ids[i] = collaborator.ID
	}

	return ids
}

// Helper function for JSON marshaling
func reviewsToIDs(reviews []*BranchReview) []uint {
	ids := make([]uint, len(reviews))

	for i, review := range reviews {
		ids[i] = review.ID
	}

	return ids
}

// Helper function for JSON marshaling
func discussionsIntoIDs(discussions []*Discussion) []uint {
	ids := make([]uint, len(discussions))

	for i, discussion := range discussions {
		ids[i] = discussion.ID
	}

	return ids
}
