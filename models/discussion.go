package models

import (
	"time"
)

type Discussion struct {
	Text      string
	Author    Member
	Deleted   bool
	Anonymous bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Replies   []Discussion
}
