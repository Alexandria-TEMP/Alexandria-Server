package tags

type PostType int16

const (
	Project PostType = iota
	Question
	Reflection
)

type PostTypeTag struct{}
