package tags

type PostReviewStatus int16

const (
	Open PostReviewStatus = iota
	RevisionNeeded
	Reviewed
)

type PostReviewStatusTag struct {
	label string
}

func (tag *PostReviewStatusTag) GetLabel() string {
	return tag.label
}

func (tag *PostReviewStatusTag) GetType() TagType {
	return PostReviewStatusType
}
