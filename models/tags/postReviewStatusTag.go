package tags

type PostReviewStatus string

const (
	Open           PostReviewStatus = "open"
	RevisionNeeded PostReviewStatus = "revision needed"
	Reviewed       PostReviewStatus = "reviewed"
)

func (tag *PostReviewStatus) GetLabel() string {
	return string(*tag)
}

func (tag *PostReviewStatus) GetType() TagType {
	return PostReviewStatusType
}
