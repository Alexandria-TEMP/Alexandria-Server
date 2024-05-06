package tags

type FeedbackPreference int16

const (
	Discussion FeedbackPreference = iota
	FormalFeedback
)

type FeedbackPreferenceTag struct{}
