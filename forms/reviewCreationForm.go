package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type ReviewCreationForm struct {
	BranchID uint

	ReviewingMemberID    uint                        `json:"reviewingMemberID"`
	BranchReviewDecision models.BranchReviewDecision `json:"branchReviewDecision"`
	Feedback             string                      `json:"feedback"`
}

// Whether the form itself contains valid data. Should NOT contain business logic (such as "if Foo > 0, Bar may not be 1")
func (form *ReviewCreationForm) IsValid() bool {
	return form.BranchReviewDecision.IsValid()
}
