package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model

	// Version has one Repository
	Repository Repository `gorm:"foreignKey:VersionID"`

	// Version has many Discussion
	Discussions []Discussion `gorm:"foreignKey:VersionID"`
}
