package tags

type PostType int16

const (
	Project PostType = iota
	Question
	Reflection
)

type PostTypeTag struct {
	label string
}

func (tag *PostTypeTag) GetLabel() string {
	return tag.label
}

func (tag *PostTypeTag) GetType() TagType {
	return PostTypeType
}
