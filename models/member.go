package models

type Member struct {
	// for lack of UUID for now
	UserID      uint64
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Institution string
	Posts       []Post
	Discussions []Discussion
	Reviews     []MergeRequestReview
}
