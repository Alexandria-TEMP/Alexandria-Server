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
	Member   Member
	MemberID uint

	// Discussion (optionally) belongs to Discussion
	Parent   *Discussion
	ParentID uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Text      string
	Deleted   bool
	Anonymous bool
}
