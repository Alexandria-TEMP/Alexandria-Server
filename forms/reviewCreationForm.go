package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

type ReviewCreationForm struct {
	// branch ID is part of endpoint

	ReviewingMemberID    uint
	BranchReviewDecision models.BranchReviewDecision
	Feedback             string
}
