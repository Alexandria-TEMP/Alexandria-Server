package tags

type ScientificField string

const (
	Mathematics     ScientificField = "mathematics"
	ComputerScience ScientificField = "computer science"
)

func (tag *ScientificField) GetLabel() string {
	return string(*tag)
}

func (tag *ScientificField) GetType() TagType {
	return ScientificFieldType
}
