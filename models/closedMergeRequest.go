package models

import "time"

type ClosedMergeRequest struct {
	CreatedAt time.Time
	MergeRequest
	MainVersionWhenClosed Version
	MergeRequestDecision
}
