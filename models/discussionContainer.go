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
	ID            uint   `json:"id"`
	DiscussionIDs []uint `json:"discussionIDs"`
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
	CurrentDiscussionContainerID     uint                               `json:"currentDiscussionContainerID"`
	MergedBranchDiscussionContainers []DiscussionContainerWithBranchDTO `json:"mergedBranchDiscussionContainers"`
}

// Represents a discussion container plus the branch it originated from
type DiscussionContainerWithBranchDTO struct {
	DiscussionContainerID uint `json:"discussionContainerID"`
	ClosedBranchID        uint `json:"closedBranchID"`
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
