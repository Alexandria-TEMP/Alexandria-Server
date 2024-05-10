package tags

type TagType int

const (
	CompletionStatusType TagType = iota
	FeedbackPreferenceType
	PostReviewStatusType
	ScientificFieldType
	PostTypeType
)

type Tag interface {
	GetLabel() string
	GetType() TagType
}
