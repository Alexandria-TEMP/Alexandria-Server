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
	ID                            uint   `json:"id" example:"1"`
	FirstName                     string `json:"firstName" example:"John"`
	LastName                      string `json:"lastName" example:"Doe"`
	Email                         string `json:"email" example:"example@example.example"`
	Institution                   string `json:"institution" example:"Example University"`
	ScientificFieldTagContainerID uint   `json:"scientificFieldTagContainerID" example:"1"`
}

type LoggedInMemberDTO struct {
	Member       MemberDTO `json:"member"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type TokenPairDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type MemberShortFormDTO struct {
	ID        uint   `json:"id" example:"1"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
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
		model.Institution,
		model.ScientificFieldTagContainerID,
	}
}

func (model *Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
