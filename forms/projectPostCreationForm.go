package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"

type ProjectPostCreationForm struct {
	PostCreationForm   PostCreationForm
	CompletionStatus   tags.CompletionStatus
	FeedbackPreference tags.FeedbackPreference
}
