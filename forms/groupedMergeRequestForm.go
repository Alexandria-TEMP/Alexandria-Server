package forms

// Holds IDs of MergeRequests
// Categorized by their MergeRequestReviewStatus
type GroupedMergeRequestForm struct {
	OpenForReviewIDs []uint
	RejectedIDs      []uint
	PeerReviewedIDs  []uint
}
