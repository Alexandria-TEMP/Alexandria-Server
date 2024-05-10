package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model

	// Version has one Repository
	Repository Repository

	// Version has many Discussion
	Discussions []Discussion
}
