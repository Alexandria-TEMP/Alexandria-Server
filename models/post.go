package models

import (
	"encoding/json"
	"slices"
	"time"

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

	Title    string
	PostType PostType
	// Post has a ScientificFieldTagContainer
	ScientificFieldTagContainer   ScientificFieldTagContainer `gorm:"foreignKey:ScientificFieldTagContainerID"`
	ScientificFieldTagContainerID uint
	// Post has a DiscussionContainer
	DiscussionContainer   DiscussionContainer `gorm:"foreignKey:DiscussionContainerID"`
	DiscussionContainerID uint

	RenderStatus RenderStatus
}

type PostDTO struct {
	ID                            uint         `json:"id" example:"1"`
	CollaboratorIDs               []uint       `json:"collaboratorIDs" example:"1"`
	Title                         string       `json:"title" example:"Post Title"`
	PostType                      PostType     `json:"postType" example:"question"`
	ScientificFieldTagContainerID uint         `json:"scientificFieldTagContainerID" example:"1"`
	DiscussionContainerID         uint         `json:"discussionContainerID" example:"1"`
	RenderStatus                  RenderStatus `json:"renderStatus" example:"success"`
	CreatedAt                     time.Time    `json:"createdAt" example:"2024-06-16T16:00:43.234Z"`
	UpdatedAt                     time.Time    `json:"updatedAt" example:"2024-06-16T16:00:43.234Z"`
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
		model.ScientificFieldTagContainerID,
		model.DiscussionContainerID,
		model.RenderStatus,
		model.CreatedAt,
		model.UpdatedAt,
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
