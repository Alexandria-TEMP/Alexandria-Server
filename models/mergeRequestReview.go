package models

import "gorm.io/gorm"

type MergeRequestDecision int16

const (
	Rejected MergeRequestDecision = iota
	Approved
)

type MergeRequestReview struct {
	gorm.Model
	Feedback string
	MergeRequestDecision
}
