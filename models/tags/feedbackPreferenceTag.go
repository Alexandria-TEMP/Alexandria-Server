package tags

type FeedbackPreference string

const (
	Discussion     FeedbackPreference = "discussion"
	FormalFeedback FeedbackPreference = "formal feedback"
)

func (tag *FeedbackPreference) GetLabel() string {
	return string(*tag)
}

func (tag *FeedbackPreference) GetType() TagType {
	return FeedbackPreferenceType
}
