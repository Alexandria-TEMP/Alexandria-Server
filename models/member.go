package models

import "gorm.io/gorm"

// A member is a logged-in user of the Alexandria platform.
type Member struct {
	gorm.Model

	// TODO can a Post be owned by multiple members? UML is not clear
	// Posts []Post

	FirstName   string
	LastName    string
	Email       string
	Password    string
	Institution string
}
