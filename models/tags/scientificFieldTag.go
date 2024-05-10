package tags

type ScientificField int16

const (
	Mathematics ScientificField = iota
)

func (tag *ScientificField) GetLabel() string {
	switch *tag {
	case Mathematics:
		return "Mathematics"
	default:
		panic("could not convert tag to string") // TODO better cleanup?
	}
}

func (tag *ScientificField) GetType() TagType {
	return ScientificFieldType
}
