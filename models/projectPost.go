package models

import (
	"encoding/json"
	"slices"

	"gorm.io/gorm"
)

// The name of the branch that will be created, automatically, when a new project post is created.
// This branch is created for purpose of peer reviewing the project post itself, before it can
// receive any other proposed changes.
const InitialPeerReviewBranchName = "Initial peer review changes"

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

// The review status of an entire Project Post
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

	CompletionStatus   ProjectCompletionStatus
	FeedbackPreference ProjectFeedbackPreference
	PostReviewStatus   ProjectReviewStatus
}

type ProjectPostDTO struct {
	ID                 uint                      `json:"id"`
	PostDTO            PostDTO                   `json:"post"`
	OpenBranchIDs      []uint                    `json:"openBranchIDs"`
	ClosedBranchIDs    []uint                    `json:"closedBranchIDs"`
	CompletionStatus   ProjectCompletionStatus   `json:"completionStatus"`
	FeedbackPreference ProjectFeedbackPreference `json:"feedbackPreference"`
	PostReviewStatus   ProjectReviewStatus       `json:"postReviewStatus"`
}

func (model *ProjectPost) GetID() uint {
	return model.Model.ID
}

func (model *ProjectPost) IntoDTO() ProjectPostDTO {
	return ProjectPostDTO{
		model.ID,
		model.Post.IntoDTO(),
		branchesToIDs(model.OpenBranches),
		closedBranchesToIDs(model.ClosedBranches),
		model.CompletionStatus,
		model.FeedbackPreference,
		model.PostReviewStatus,
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
