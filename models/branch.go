package models

import (
	"encoding/json"
	"slices"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type RenderStatus string

const (
	Success RenderStatus = "success"
	Pending RenderStatus = "pending"
	Failure RenderStatus = "failure"
)

func (enum *RenderStatus) IsValid() bool {
	valid := []RenderStatus{Success, Pending, Failure}
	return slices.Contains(valid, *enum)
}

// The branch's aggregated review status, derived from its individual reviews' statuses
type BranchOverallReviewStatus string

const (
	BranchOpenForReview BranchOverallReviewStatus = "open for review"
	BranchPeerReviewed  BranchOverallReviewStatus = "peer reviewed"
	BranchRejected      BranchOverallReviewStatus = "rejected"
)

func (enum *BranchOverallReviewStatus) IsValid() bool {
	valid := []BranchOverallReviewStatus{BranchOpenForReview, BranchPeerReviewed, BranchRejected}
	return slices.Contains(valid, *enum)
}

type Branch struct {
	gorm.Model

	/////////////////////////////////////////////
	// The branch's proposed changes:

	NewPostTitle string

	UpdatedCompletionStatus ProjectCompletionStatus
	UpdatedScientificFields []tags.ScientificField `gorm:"serializer:json"`

	/////////////////////////////////////////////
	// The branch's metadata:

	// Branch has many BranchCollaborator
	Collaborators []*BranchCollaborator `gorm:"foreignKey:BranchID"`

	// Branch has many BranchReview
	Reviews []*BranchReview `gorm:"foreignKey:BranchID"`

	// Branch has a DiscussionContainer
	DiscussionContainer   DiscussionContainer `gorm:"foreignKey:DiscussionContainerID"`
	DiscussionContainerID uint

	// ProjectPost has many Branch
	ProjectPostID uint

	BranchTitle string

	RenderStatus       RenderStatus
	BranchReviewStatus BranchOverallReviewStatus
}

type BranchDTO struct {
	ID uint
	// MR's proposed changes
	NewPostTitle            string
	UpdatedCompletionStatus ProjectCompletionStatus
	UpdatedScientificFields []tags.ScientificField
	// MR metadata
	CollaboratorIDs    []uint
	ReviewIDs          []uint
	ProjectPostID      uint
	BranchTitle        string
	RenderStatus       RenderStatus
	DiscussionIDs      []uint
	BranchReviewStatus BranchOverallReviewStatus
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
		model.RenderStatus,
		discussionContainerIntoIDs(&model.DiscussionContainer),
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
func discussionContainerIntoIDs(discussions *DiscussionContainer) []uint {
	ids := make([]uint, len(discussions.Discussions))

	for i, discussion := range discussions.Discussions {
		ids[i] = discussion.ID
	}

	return ids
}
