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

type ReviewStatus string

const (
	BranchOpenForReview ReviewStatus = "open for review"
	BranchPeerReviewed  ReviewStatus = "peer reviewed"
	BranchRejected      ReviewStatus = "rejected"
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

	// Branch has many Review
	Reviews []*Review `gorm:"foreignKey:BranchID"`

	// Branch has a DiscussionContainer
	DiscussionContainer   DiscussionContainer `gorm:"foreignKey:DiscussionContainerID"`
	DiscussionContainerID uint

	// ProjectPost has many Branch
	ProjectPostID uint

	BranchTitle string

	Anonymous bool

	RenderStatus RenderStatus
	ReviewStatus ReviewStatus
}

type BranchDTO struct {
	ID uint `json:"id"`
	// Branch's proposed changes
	NewPostTitle            string                 `json:"new_post_title"`
	UpdatedCompletionStatus tags.CompletionStatus  `json:"updated_completion_status"`
	UpdatedScientificFields []tags.ScientificField `json:"updated_scientific_fields"`
	// Branch metadata
	CollaboratorIDs []uint       `json:"collaborator_ids"`
	ReviewIDs       []uint       `json:"review_ids"`
	ProjectPostID   uint         `json:"project_post_id"`
	BranchTitle     string       `json:"branch_title"`
	Anonymous       bool         `json:"anonymous"`
	RenderStatus    RenderStatus `json:"render_status"`
	DiscussionIDs   []uint       `json:"discussion_ids"`
	ReviewStatus    ReviewStatus `json:"review_status"`
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
		discussionContainerIntoIDs(&model.DiscussionContainer),
		model.ReviewStatus,
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
func reviewsToIDs(reviews []*Review) []uint {
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
