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
	// Member has a ScientificFieldTagContainer
	ScientificFieldTagContainer   ScientificFieldTagContainer `gorm:"foreignKey:ScientificFieldTagContainerID"`
	ScientificFieldTagContainerID uint
}

type MemberDTO struct {
	ID                            uint   `json:"id"`
	FirstName                     string `json:"firstName"`
	LastName                      string `json:"lastName"`
	Email                         string `json:"email"`
	Password                      string `json:"password"`
	Institution                   string `json:"institution"`
	ScientificFieldTagContainerID uint   `json:"scientificFieldTagContainerID"`
}

type MemberShortFormDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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
		model.ScientificFieldTagContainerID,
	}
}

func (model *Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
