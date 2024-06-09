package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type DiscussionService struct {
	DiscussionRepository database.ModelRepositoryInterface[*models.Discussion]
}

func (discussionService *DiscussionService) GetDiscussion(id uint) (*models.Discussion, error) {
	return discussionService.DiscussionRepository.GetByID(id)
}

func (discussionService *DiscussionService) CreateRootDiscussion(_ *forms.RootDiscussionCreationForm) (*models.Discussion, error) {
	// TODO implement
	return nil, fmt.Errorf("todo")
}

func (discussionService *DiscussionService) CreateReply(_ *forms.ReplyDiscussionCreationForm) (*models.Discussion, error) {
	// TODO implement
	return nil, fmt.Errorf("todo")
}
