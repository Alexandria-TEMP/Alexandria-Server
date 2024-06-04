package forms

// Holds IDs of Branches
// Categorized by their BranchReviewStatus
type GroupedBranchForm struct {
	OpenForReviewIDs []uint
	RejectedIDs      []uint
	PeerReviewedIDs  []uint
}
