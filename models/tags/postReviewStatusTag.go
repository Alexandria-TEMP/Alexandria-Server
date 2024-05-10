package tags

type PostReviewStatus int16

const (
	Open PostReviewStatus = iota
	RevisionNeeded
	Reviewed
)

func (tag *PostReviewStatus) GetLabel() string {
	switch *tag {
	case Open:
		return "Open"
	case RevisionNeeded:
		return "Revision Needed"
	case Reviewed:
		return "Reviewed"
	default:
		panic("could not convert tag to string") // TODO better cleanup?
	}
}

func (tag *PostReviewStatus) GetType() TagType {
	return PostReviewStatusType
}
