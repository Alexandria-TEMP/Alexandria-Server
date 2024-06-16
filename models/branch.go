package models

import (
	"encoding/json"
	"slices"
	"time"

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

// The branch's aggregated branchreview status, derived from its individual reviews' statuses
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

	UpdatedPostTitle                     *string
	UpdatedCompletionStatus              *ProjectCompletionStatus
	UpdatedFeedbackPreferences           *ProjectFeedbackPreference
	UpdatedScientificFieldTagContainer   *ScientificFieldTagContainer `gorm:"foreignKey:UpdatedScientificFieldTagContainerID"`
	UpdatedScientificFieldTagContainerID *uint

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
	ProjectPostID *uint

	BranchTitle string

	RenderStatus              RenderStatus
	BranchOverallReviewStatus BranchOverallReviewStatus
}

type BranchDTO struct {
	ID uint `json:"id"`
	// MR's proposed changes
	UpdatedPostTitle                     *string                  `json:"UpdatedPostTitle"`
	UpdatedCompletionStatus              *ProjectCompletionStatus `json:"updatedCompletionStatus"`
	UpdatedScientificFieldTagContainerID *uint                    `json:"updatedScientificFieldTagContainerID"`
	// MR metadata
	CollaboratorIDs           []uint                    `json:"collaboratorIDs"`
	ReviewIDs                 []uint                    `json:"reviewIDs"`
	ProjectPostID             *uint                     `json:"projectPostID"`
	BranchTitle               string                    `json:"branchTitle"`
	RenderStatus              RenderStatus              `json:"renderStatus"`
	DiscussionContainerID     uint                      `json:"discussionContainerID"`
	BranchOverallReviewStatus BranchOverallReviewStatus `json:"branchOverallReviewStatus"`
	CreatedAt                 time.Time                 `json:"createdAt"`
	UpdatedAt                 time.Time                 `json:"updatedAt"`
}

func (model *Branch) GetID() uint {
	return model.Model.ID
}

func (model *Branch) IntoDTO() BranchDTO {
	return BranchDTO{
		model.ID,
		model.UpdatedPostTitle,
		model.UpdatedCompletionStatus,
		model.UpdatedScientificFieldTagContainerID,
		branchCollaboratorsToIDs(model.Collaborators),
		reviewsToIDs(model.Reviews),
		model.ProjectPostID,
		model.BranchTitle,
		model.RenderStatus,
		model.DiscussionContainerID,
		model.BranchOverallReviewStatus,
		model.CreatedAt,
		model.UpdatedAt,
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

	for i, branchreview := range reviews {
		ids[i] = branchreview.ID
	}

	return ids
}

// Holds IDs of Branches and ClosedBranches
// Categorized by their BranchReviewStatus
type BranchesGroupedByReviewStatusDTO struct {
	OpenBranchIDs           []uint `json:"openBranchIDs"`
	RejectedClosedBranchIDs []uint `json:"rejectedClosedBranchIDs"`
	ApprovedClosedBranchIDs []uint `json:"approvedClosedBranchIDs"`
}
