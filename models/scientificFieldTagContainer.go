package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ScientificFieldTagContainer struct {
	gorm.Model

	// ScientificFieldTagContainer has many ScientificFieldTag
	ScientificFieldTags []*ScientificFieldTag `gorm:"many2many:tag_containers;"`
}

type ScientificFieldTagContainerDTO struct {
	ID                    uint   `json:"id" example:"1"`
	ScientificFieldTagIDs []uint `json:"scientificFieldTagIDs" example:"1"`
}

func (model *ScientificFieldTagContainer) GetID() uint {
	return model.Model.ID
}

func (model *ScientificFieldTagContainer) IntoDTO() ScientificFieldTagContainerDTO {
	return ScientificFieldTagContainerDTO{
		ID:                    model.ID,
		ScientificFieldTagIDs: ScientificFieldTagIntoIDs(model.ScientificFieldTags),
	}
}

func (model *ScientificFieldTagContainer) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
