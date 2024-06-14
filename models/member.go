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
	ID                    uint   `json:"id"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Email                 string `json:"email"`
	Institution           string `json:"institution"`
	ScientificFieldTagIDs []uint `json:"scientificFieldTagIDs"`
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
		model.Institution,
		ScientificFieldTagContainerIntoIDs(&model.ScientificFieldTagContainer),
	}
}

func (model *Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
