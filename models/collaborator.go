package models

type CollaborationType int16

const (
	Author CollaborationType = iota
	Contributor
	Reviewer
)

type Collaborator struct {
	collaboratorID uint64
	Member
	CollaborationType
}
