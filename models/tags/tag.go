package tags

type TagType int16

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
