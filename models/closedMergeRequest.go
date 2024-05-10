package models

import (
	"time"

	"gorm.io/gorm"
)

type ClosedMergeRequest struct {
	gorm.Model
	CreatedAt time.Time
	MergeRequest
	MainVersionWhenClosed Version
	MergeRequestDecision
}
