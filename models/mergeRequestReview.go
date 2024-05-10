package models

import "gorm.io/gorm"

type MergeRequestDecision int16

const (
	Rejected MergeRequestDecision = iota
	Approved
)

type MergeRequestReview struct {
	gorm.Model

	// MergeRequest has many MergeRequestReview
	MergeRequestID uint

	// MergeRequestReview belongs to Member
	Member   Member
	MemberID uint

	MergeRequestDecision
	Feedback string
}
