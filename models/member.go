package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// A member is a logged-in user of the Alexandria platform.
type Member struct {
	gorm.Model

	FirstName        string
	LastName         string
	Email            string
	Password         string // TODO hmmmmmm maybe not
	Institution      string
	ScientificFields []ScientificField `gorm:"serializer:json"`
}

type MemberDTO struct {
	ID               uint              `json:"id"`
	FirstName        string            `json:"firstName"`
	LastName         string            `json:"lastName"`
	Email            string            `json:"email"`
	Password         string            `json:"password"`
	Institution      string            `json:"institution"`
	ScientificFields []ScientificField `json:"ScientificFields"`
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
		model.ScientificFields,
	}
}

func (model *Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
