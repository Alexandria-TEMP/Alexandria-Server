package forms

// Holds IDs of Branches
// Categorized by their ReviewStatus
type GroupedBranchForm struct {
	OpenForReviewIDs []uint
	RejectedIDs      []uint
	PeerReviewedIDs  []uint
}
