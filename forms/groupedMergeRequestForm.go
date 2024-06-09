package forms

// Holds IDs of Branches
// Categorized by their BranchReviewStatus
type GroupedBranchForm struct {
	OpenForReviewIDs []uint
	RejectedIDs      []uint
	PeerReviewedIDs  []uint
}

// TODO this should not be a form!

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *GroupedBranchForm) IsValid() bool {
	return true
}
