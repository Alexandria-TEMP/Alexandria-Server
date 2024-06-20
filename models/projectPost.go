package models

import (
	"encoding/json"
	"slices"
	"time"

	"gorm.io/gorm"
)

// The name of the branch that will be created, automatically, when a new project post is created.
// This branch is created for purpose of peer reviewing the project post itself, before it can
// receive any other proposed changes.
const InitialPeerReviewBranchName = "Initial Peer Review"

type ProjectCompletionStatus string

const (
	Idea      ProjectCompletionStatus = "idea"
	Ongoing   ProjectCompletionStatus = "ongoing"
	Completed ProjectCompletionStatus = "completed"
)

func (enum *ProjectCompletionStatus) IsValid() bool {
	valid := []ProjectCompletionStatus{Idea, Ongoing, Completed}
	return slices.Contains(valid, *enum)
}

type ProjectFeedbackPreference string

const (
	DiscussionFeedback ProjectFeedbackPreference = "discussion feedback"
	FormalFeedback     ProjectFeedbackPreference = "formal feedback"
)

func (enum *ProjectFeedbackPreference) IsValid() bool {
	valid := []ProjectFeedbackPreference{DiscussionFeedback, FormalFeedback}
	return slices.Contains(valid, *enum)
}

// The branchreview status of an entire Project Post
// If a Project Post is not (yet) peer reviewed, new changes cannot be requested
type ProjectReviewStatus string

const (
	Open           ProjectReviewStatus = "open"
	RevisionNeeded ProjectReviewStatus = "revision needed"
	Reviewed       ProjectReviewStatus = "reviewed"
)

func (enum *ProjectReviewStatus) IsValid() bool {
	valid := []ProjectReviewStatus{Open, RevisionNeeded, Reviewed}
	return slices.Contains(valid, *enum)
}

type ProjectPost struct {
	gorm.Model

	// ProjectPost belongs to Post
	Post   Post `gorm:"foreignKey:PostID"`
	PostID uint

	// ProjectPost has many Branch
	OpenBranches []*Branch `gorm:"foreignKey:ProjectPostID"`

	// ProjectPost has many ClosedBranch
	ClosedBranches []*ClosedBranch `gorm:"foreignKey:ProjectPostID"`

	ProjectCompletionStatus   ProjectCompletionStatus
	ProjectFeedbackPreference ProjectFeedbackPreference
	PostReviewStatus          ProjectReviewStatus
}

type ProjectPostDTO struct {
	ID                        uint                      `json:"id" example:"1"`
	PostID                    uint                      `json:"postID" example:"1"`
	OpenBranchIDs             []uint                    `json:"openBranchIDs" example:"1"`
	ClosedBranchIDs           []uint                    `json:"closedBranchIDs" example:"1"`
	ProjectCompletionStatus   ProjectCompletionStatus   `json:"projectCompletionStatus" example:"ongoing"`
	ProjectFeedbackPreference ProjectFeedbackPreference `json:"projectFeedbackPreference" example:"formal feedback"`
	PostReviewStatus          ProjectReviewStatus       `json:"postReviewStatus" example:"open"`
	CreatedAt                 time.Time                 `json:"createdAt" example:"2024-06-16T16:00:43.234Z"`
	UpdatedAt                 time.Time                 `json:"updatedAt" example:"2024-06-16T16:00:43.234Z"`
}

func (model *ProjectPost) GetID() uint {
	return model.Model.ID
}

func (model *ProjectPost) IntoDTO() ProjectPostDTO {
	return ProjectPostDTO{
		model.ID,
		model.PostID,
		branchesToIDs(model.OpenBranches),
		closedBranchesToIDs(model.ClosedBranches),
		model.ProjectCompletionStatus,
		model.ProjectFeedbackPreference,
		model.PostReviewStatus,
		model.CreatedAt,
		model.UpdatedAt,
	}
}

func (model *ProjectPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func branchesToIDs(branches []*Branch) []uint {
	ids := make([]uint, len(branches))

	for i, branches := range branches {
		ids[i] = branches.ID
	}

	return ids
}

// Helper function for JSON marshaling
func closedBranchesToIDs(branches []*ClosedBranch) []uint {
	ids := make([]uint, len(branches))

	for i, branches := range branches {
		ids[i] = branches.ID
	}

	return ids
}
