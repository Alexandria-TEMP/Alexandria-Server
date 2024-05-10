package models

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Institution string
	Posts       []Post
	Discussions []Discussion
	Reviews     []MergeRequestReview
}
