package models

import "gorm.io/gorm"

// A member is a logged-in user of the Alexandria platform.
type Member struct {
	gorm.Model

	FirstName   string
	LastName    string
	Email       string
	Password    string
	Institution string
}

func (model *Member) GetID() uint {
	return model.Model.ID
}
