package tags

type CompletionStatus int16

const (
	Idea CompletionStatus = iota
	Ongoing
	Completed
)

type CompletionStatusTag struct {
	label string
}

func (tag *CompletionStatusTag) GetLabel() string {
	return tag.label
}

func (tag *CompletionStatusTag) GetType() TagType {
	return CompletionStatusType
}
