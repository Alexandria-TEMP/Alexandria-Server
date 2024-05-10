package models

import (
	"time"

	"gorm.io/gorm"
)

type ClosedMergeRequest struct {
	gorm.Model

	// ClosedMergeRequest belongs to MergeRequest
	MergeRequest   MergeRequest
	MergeRequestID uint

	// ClosedMergeRequest belongs to Version
	MainVersionWhenClosed   Version
	MainVersionWhenClosedID uint

	// ProjectPost has many ClosedMergeRequest
	ProjectPostID uint

	CreatedAt            time.Time
	MergeRequestDecision MergeRequestDecision
}
