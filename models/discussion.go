package models

import (
	"time"

	"gorm.io/gorm"
)

type Discussion struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Text      string
	Author    Member
	Replies   []Discussion
	Deleted   bool
	Anonymous bool
}
