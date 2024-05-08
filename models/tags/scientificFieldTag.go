package tags

type ScientificField int16

const (
	Maf ScientificField = iota
)

type ScientificFieldTag struct {
	label string
}

func (tag *ScientificFieldTag) GetLabel() string {
	return tag.label
}

func (tag *ScientificFieldTag) GetType() TagType {
	return ScientificFieldType
}
