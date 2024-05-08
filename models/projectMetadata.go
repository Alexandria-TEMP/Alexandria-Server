package models

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"

type ProjectMetadata struct {
	CompletionStatus    tags.CompletionStatusTag
	FeedbackPreference  tags.FeedbackPreferenceTag
	PostReviewStatusTag tags.PostReviewStatusTag
	ForkedFrom          ClosedMergeRequest
}
