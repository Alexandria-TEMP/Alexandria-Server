package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model

	// Version has many Discussion
	Discussions []Discussion `gorm:"foreignKey:VersionID"`
}
