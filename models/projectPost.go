package models

type ProjectPost struct {
	ID uint64
	Post
	ProjectMetadata
	OpenMergeRequests   []MergeRequest
	ClosedMergeRequests []ClosedMergeRequest
}
