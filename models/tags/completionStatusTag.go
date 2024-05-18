package tags

type CompletionStatus string

const (
	Idea      CompletionStatus = "idea"
	Ongoing   CompletionStatus = "ongoing"
	Completed CompletionStatus = "completed"
)

func (tag *CompletionStatus) GetLabel() string {
	return string(*tag)
}

func (tag *CompletionStatus) GetType() TagType {
	return CompletionStatusType
}
