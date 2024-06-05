package tags

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ScientificFieldTagContainer struct {
	gorm.Model

	// ScientificFieldTagContainer has many ScientificFieldTag
	ScientificFieldTags []*ScientificFieldTag `gorm:"foreignKey:ContainerID"`
}

// a scientific field tag is a tag representing a specific scientific field
type ScientificFieldTag struct {
	gorm.Model

	ScientificField string
	// ScientificFieldTag belongs to ScientificFieldTagContainer
	ContainerID uint
	// Tag can optionally have many subtags, or many ScientificFieldTag
	Subtags  []*ScientificFieldTag `gorm:"foreignKey:ParentID"`
	ParentID *uint
}

type ScientificFieldTagDTO struct {
	ID              uint
	ScientificField string
	SubtagIDs       []uint
}

func (model *ScientificFieldTag) GetID() uint {
	return model.Model.ID
}

func (model *ScientificFieldTag) IntoDTO() ScientificFieldTagDTO {
	return ScientificFieldTagDTO{
		model.ID,
		model.ScientificField,
		ScientificFieldTagIntoIDs(model.Subtags),
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

// Helper function for JSON marshaling
func ScientificFieldTagContainerIntoIDs(scientificFieldTags *ScientificFieldTagContainer) []uint {
	ids := make([]uint, len(scientificFieldTags.ScientificFieldTags))

	for i, tag := range scientificFieldTags.ScientificFieldTags {
		ids[i] = tag.ID
	}

	return ids
}
