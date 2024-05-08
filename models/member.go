package models

type Member struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Institution string
	Posts       []Post
	Discussions []Discussion
	Reviews     []MergeRequestReview
}
