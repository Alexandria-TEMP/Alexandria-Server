package models

import (
	"encoding/json"
	"slices"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type PostType string

const (
	Project    PostType = "project"
	Question   PostType = "question"
	Reflection PostType = "reflection"
)

func (enum *PostType) IsValid() bool {
	valid := []PostType{Project, Question, Reflection}
	return slices.Contains(valid, *enum)
}

type Post struct {
	gorm.Model

	// Post has many PostCollaborator
	Collaborators []*PostCollaborator `gorm:"foreignKey:PostID"`

	// Post files and render can be implicitly accessed in the vfs with the postID

	Title               string
	PostType            PostType
	ScientificFieldTags []tags.ScientificField `gorm:"serializer:json"`

	// Post has a DiscussionContainer
	DiscussionContainer   DiscussionContainer `gorm:"foreignKey:DiscussionContainerID"`
	DiscussionContainerID uint
}

type PostDTO struct {
	ID                  uint                   `json:"id"`
	CollaboratorIDs     []uint                 `json:"collaboratorIDs"`
	Title               string                 `json:"title"`
	PostType            PostType               `json:"postType"`
	ScientificFieldTags []tags.ScientificField `json:"scientificFieldTags"`
	DiscussionIDs       []uint                 `json:"discussionIDs"`
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
		model.ScientificFieldTags,
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
