package tags

type FeedbackPreference int16

const (
	Discussion FeedbackPreference = iota
	FormalFeedback
)

func (tag *FeedbackPreference) GetLabel() string {
	switch *tag {
	case Discussion:
		return "Discussion"
	case FormalFeedback:
		return "Formal Feedback"
	default:
		panic("could not convert tag to string") // TODO better cleanup?
	}
}

func (tag *FeedbackPreference) GetType() TagType {
	return FeedbackPreferenceType
}
