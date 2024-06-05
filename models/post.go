package models

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	// Post has many PostCollaborator
	Collaborators []*PostCollaborator `gorm:"foreignKey:PostID"`

	// Post files and render can be implicitly accessed in the vfs with the postID

	Title               string
	PostType            tags.PostType
	ScientificFieldTags []*tags.ScientificFieldTag `gorm:"foreignKey:PostID"`
	// Post has a DiscussionContainer
	DiscussionContainer   DiscussionContainer `gorm:"foreignKey:DiscussionContainerID"`
	DiscussionContainerID uint
}

type PostDTO struct {
	ID                    uint
	CollaboratorIDs       []uint
	Title                 string
	PostType              tags.PostType
	ScientificFieldTagIDs []uint
	DiscussionIDs         []uint
}

func (model *Post) GetID() uint {
	return model.Model.ID
}

func (model *Post) IntoDTO() PostDTO {
	return PostDTO{
		model.ID,
		postCollaboratorsToIDs(model.Collaborators),
		model.Title,
		model.PostType,
		tags.ScientificFieldTagIntoIDs(model.ScientificFieldTags),
		discussionContainerIntoIDs(&model.DiscussionContainer),
	}
}

func (model *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func postCollaboratorsToIDs(collaborators []*PostCollaborator) []uint {
	ids := make([]uint, len(collaborators))

	for i, collaborator := range collaborators {
		ids[i] = collaborator.ID
	}

	return ids
}
