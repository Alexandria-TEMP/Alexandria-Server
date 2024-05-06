package models

type MergeRequestDecision int16

const (
	Rejected MergeRequestDecision = iota
	Approved
)

type MergeRequestReview struct{}
