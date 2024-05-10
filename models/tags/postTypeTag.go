package tags

type PostType int16

const (
	Project PostType = iota
	Question
	Reflection
)

func (tag *PostType) GetLabel() string {
	switch *tag {
	case Project:
		return "Project"
	case Question:
		return "Question"
	case Reflection:
		return "Reflection"
	default:
		panic("could not convert tag to string") // TODO better cleanup?
	}
}

func (tag *PostType) GetType() TagType {
	return PostTypeType
}
