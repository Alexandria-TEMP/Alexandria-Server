package models

import "time"

type ClosedMergeRequest struct {
	MergeRequest
	MergeRequestDecision
	MainVersionWhenClosed Version
	CreatedAt             time.Time
}
