package forms

// Holds IDs of Branches
// Categorized by their BranchOverallReviewStatus
type GroupedBranchForm struct {
	OpenForReviewIDs []uint `json:"openForReviewIDs"`
	RejectedIDs      []uint `json:"rejectedIDs"`
	PeerReviewedIDs  []uint `json:"peerReviewedIDs"`
}

// TODO this should not be a form!

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *GroupedBranchForm) IsValid() bool {
	return true
}
