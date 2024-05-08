package models

type ProjectPost struct {
	Post
	ProjectMetadata
	OpenMergeRequests   []MergeRequest
	ClosedMergeRequests []ClosedMergeRequest
}
