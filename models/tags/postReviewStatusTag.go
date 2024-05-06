package tags

type PostReviewStatus int16

const (
	Open PostReviewStatus = iota
	RevisionNeeded
	Reviewed
)

type PostReviewStatusTag struct{}
