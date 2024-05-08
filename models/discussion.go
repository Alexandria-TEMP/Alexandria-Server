package models

import (
	"time"
)

type Discussion struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Text      string
	Author    Member
	Replies   []Discussion
	Deleted   bool
	Anonymous bool
}
