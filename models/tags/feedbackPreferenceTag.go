package tags

type FeedbackPreference int16

const (
	Discussion FeedbackPreference = iota
	FormalFeedback
)

type FeedbackPreferenceTag struct {
	label string
}

func (tag *FeedbackPreferenceTag) GetLabel() string {
	return tag.label
}

func (tag *FeedbackPreferenceTag) GetType() TagType {
	return FeedbackPreferenceType
}
