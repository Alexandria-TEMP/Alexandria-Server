package models

import (
	"time"

	"gorm.io/gorm"
)

type ClosedMergeRequest struct {
	gorm.Model

	// ClosedMergeRequest belongs to MergeRequest
	MergeRequest   MergeRequest `gorm:"foreignKey:MergeRequestID"`
	MergeRequestID uint

	// ClosedMergeRequest belongs to Version
	MainVersionWhenClosed   Version `gorm:"foreignKey:MainVersionWhenClosedID"`
	MainVersionWhenClosedID uint

	// ProjectPost has many ClosedMergeRequest
	ProjectPostID uint

	CreatedAt            time.Time
	MergeRequestDecision MergeRequestDecision `gorm:"serializer:json"`
}
