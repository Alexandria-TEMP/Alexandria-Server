package forms

type DiscussionCreationForm struct {
	// Parent ID / Version ID is sent as an optional query parameter

	// If anonymous, the discussion will ignore member ID
	Anonymous bool
	MemberID  uint

	Text string
}
