package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type DiscussionService interface {
	GetDiscussion(id uint) (*models.Discussion, error)

	// Add a discussion to a discussion container (on a Post, Branch...)
	CreateRootDiscussion(form *forms.RootDiscussionCreationForm) (*models.Discussion, error)

	// Add a discussion as a reply to another discussion
	CreateReply(form *forms.ReplyDiscussionCreationForm) (*models.Discussion, error)
}
