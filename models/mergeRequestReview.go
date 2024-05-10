package models

import "gorm.io/gorm"

type MergeRequestDecision string

const (
	Rejected MergeRequestDecision = "rejected"
	Approved MergeRequestDecision = "approved"
)

type MergeRequestReview struct {
	gorm.Model

	// MergeRequest has many MergeRequestReview
	MergeRequestID uint

	// MergeRequestReview belongs to Member
	Member   Member `gorm:"foreignKey:MemberID"`
	MemberID uint

	MergeRequestDecision MergeRequestDecision `gorm:"serializer:json"`
	Feedback             string
}
