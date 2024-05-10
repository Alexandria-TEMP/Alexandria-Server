package tags

type CompletionStatus int16

const (
	Idea CompletionStatus = iota
	Ongoing
	Completed
)

func (tag *CompletionStatus) GetLabel() string {
	switch *tag {
	case Idea:
		return "Idea"
	case Ongoing:
		return "Ongoing"
	case Completed:
		return "Completed"
	default:
		panic("could not convert tag to string") // TODO better cleanup?
	}
}

func (tag *CompletionStatus) GetType() TagType {
	return CompletionStatusType
}
