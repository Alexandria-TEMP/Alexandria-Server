package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// a scientific field tag is a tag representing a specific scientific field
type ScientificFieldTag struct {
	gorm.Model

	ScientificField string
	// ScientificFieldTag belongs to ScientificFieldTagContainer
	Containers []*ScientificFieldTagContainer `gorm:"many2many:tag_containers;"`
	// Tag can optionally have many subtags, or many ScientificFieldTag
	Subtags  []*ScientificFieldTag `gorm:"foreignKey:ParentID"`
	ParentID *uint
}

type ScientificFieldTagDTO struct {
	ID              uint   `json:"id"`
	ScientificField string `json:"scientificField"`
	SubtagIDs       []uint `json:"subtagIDs"`
	ParentID        *uint  `json:"parentID"`
}

func (model *ScientificFieldTag) GetID() uint {
	return model.Model.ID
}

func (model *ScientificFieldTag) IntoDTO() ScientificFieldTagDTO {
	return ScientificFieldTagDTO{
		model.ID,
		model.ScientificField,
		ScientificFieldTagIntoIDs(model.Subtags),
		model.ParentID,
	}
}

func (model *ScientificFieldTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}

// Helper function for JSON marshaling
func ScientificFieldTagIntoIDs(subtags []*ScientificFieldTag) []uint {
	ids := make([]uint, len(subtags))

	for i, subtag := range subtags {
		ids[i] = subtag.ID
	}

	return ids
}
