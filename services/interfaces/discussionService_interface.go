package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./discussionService_interface.go -destination=../../mocks/discussionService_mock.go

type DiscussionService interface {
	// GetDiscussion fetches a discussion from the database by its ID
	GetDiscussion(id uint) (*models.Discussion, error)

	// CreateRootDiscussion adds a discussion to a discussion container of a post, branch, etc.
	CreateRootDiscussion(form *forms.RootDiscussionCreationForm) (*models.Discussion, error)

	// CreateReply adds a discussion as a reply to another discussion
	CreateReply(form *forms.ReplyDiscussionCreationForm) (*models.Discussion, error)
}
