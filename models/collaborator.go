package models

import "gorm.io/gorm"

type CollaborationType int16

const (
	Author CollaborationType = iota
	Contributor
	Reviewer
)

type Collaborator struct {
	gorm.Model
	Member
	CollaborationType
}
