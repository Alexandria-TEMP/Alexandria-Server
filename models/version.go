package models

import (
	"gorm.io/gorm"
)

type Version struct {
	gorm.Model

	Repository Repository `gorm:"serializer:json"`

	// Version has many Discussion
	Discussions []Discussion `gorm:"foreignKey:VersionID"`
}

func (model *Version) GetID() uint {
	return model.Model.ID
}

type Repository struct {
	// TODO write serialization/deserialization, OR use a filesystem instead
	// QuartoProject multipart.File `swaggerignore:"true"`
}
