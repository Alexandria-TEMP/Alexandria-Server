package models

import (
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

type PostMetadata struct {
	gorm.Model
	Collaborators       []Collaborator
	CreatedAt           time.Time
	UpdatedAt           time.Time
	PostType            tags.PostTypeTag
	ScientificFieldTags []tags.ScientificFieldTag
}
