package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type DiscussionContainer struct {
	gorm.Model

	// DiscussionContainer has many Discussion
	Discussions []*Discussion `gorm:"foreignKey:ContainerID"`
}

type DiscussionContainerDTO struct {
	ID            uint   `json:"id" example:"1"`
	DiscussionIDs []uint `json:"discussionIDs" example:"1"`
}

func (model *DiscussionContainer) GetID() uint {
	return model.Model.ID
}

func (model *DiscussionContainer) IntoDTO() DiscussionContainerDTO {
	return DiscussionContainerDTO{
		ID:            model.ID,
		DiscussionIDs: onlyRootDiscussionsIntoIDs(model.Discussions),
	}
}

func (model *DiscussionContainer) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Links discussion containers with their place of origin, for purpose
// of sending a project post's entire discussion history.
type DiscussionContainerProjectHistoryDTO struct {
	CurrentDiscussionContainerID     uint                               `json:"currentDiscussionContainerID" example:"1"`
	MergedBranchDiscussionContainers []DiscussionContainerWithBranchDTO `json:"mergedBranchDiscussionContainers"`
}

// Represents a discussion container plus the branch it originated from
type DiscussionContainerWithBranchDTO struct {
	DiscussionContainerID uint `json:"discussionContainerID" example:"2"`
	ClosedBranchID        uint `json:"closedBranchID" example:"1"`
}

// onlyRootDiscussionsIntoIDs takes a list of discussions, and returns the IDs of all root discussions
func onlyRootDiscussionsIntoIDs(discussions []*Discussion) []uint {
	rootDiscussionIDs := []uint{}

	for _, discussion := range discussions {
		if discussion.ParentID == nil {
			rootDiscussionIDs = append(rootDiscussionIDs, discussion.ID)
		}
	}

	return rootDiscussionIDs
}
