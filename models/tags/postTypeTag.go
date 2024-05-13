package tags

type PostType string

const (
	Project    PostType = "project"
	Question   PostType = "question"
	Reflection PostType = "reflection"
)

func (tag *PostType) GetLabel() string {
	return string(*tag)
}

func (tag *PostType) GetType() TagType {
	return PostTypeType
}
