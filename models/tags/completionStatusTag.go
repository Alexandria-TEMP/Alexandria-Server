package tags

type CompletionStatus int16

const (
	Idea CompletionStatus = iota
	Ongoing
	Completed
)

type CompletionStatusTag struct{}
