package models

type Post struct {
	ID uint64 // place holder until we make standardzied uuid system
	PostMetadata
	CurrentVersion Version
}
