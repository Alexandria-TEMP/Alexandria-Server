package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type DiscussionContainerService struct {
	DiscussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
}

func (discussionContainerService *DiscussionContainerService) GetDiscussionContainer(id uint) (*models.DiscussionContainer, error) {
	return discussionContainerService.DiscussionContainerRepository.GetByID(id)
}
