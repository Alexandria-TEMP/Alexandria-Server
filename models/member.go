package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// A member is a logged-in user of the Alexandria platform.
type Member struct {
	gorm.Model

	FirstName   string
	LastName    string
	Email       string
	Password    string // TODO hmmmmmm maybe not
	Institution string
}

type MemberDTO struct {
	ID          uint
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Institution string
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
	}
}

func (model *Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
