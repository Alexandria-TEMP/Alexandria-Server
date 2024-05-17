package models

import (
	"time"

	"gorm.io/gorm"
)

type Discussion struct {
	gorm.Model

	// Version has many Discussion
	VersionID uint

	// Discussion belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	// Discussion optionally has many Discussion
	Replies  []Discussion `gorm:"foreignKey:ParentID"`
	ParentID *uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Text      string
	Deleted   bool
	Anonymous bool
}

func (model *Discussion) GetID() uint {
	return model.Model.ID
}
