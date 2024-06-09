package models

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

// A member is a logged-in user of the Alexandria platform.
type Member struct {
	gorm.Model

	FirstName           string
	LastName            string
	Email               string
	Password            string // TODO hmmmmmm maybe not
	Institution         string
	ScientificFieldTags []tags.ScientificField `gorm:"serializer:json"`
}

type MemberDTO struct {
	ID                  uint                   `json:"id"`
	FirstName           string                 `json:"firstName"`
	LastName            string                 `json:"lastName"`
	Email               string                 `json:"email"`
	Password            string                 `json:"password"`
	Institution         string                 `json:"institution"`
	ScientificFieldTags []tags.ScientificField `json:"scientificFieldTags"`
}

func (model *Member) GetID() uint {
	return model.Model.ID
}

func (model *Member) IntoDTO() MemberDTO {
	return MemberDTO{
		model.ID,
		model.FirstName,
		model.LastName,
		model.Email,
		model.Password,
		model.Institution,
		model.ScientificFieldTags,
	}
}

func (model *Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
