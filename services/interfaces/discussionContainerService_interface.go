package interfaces

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

//go:generate mockgen -package=mocks -source=./discussionContainerService_interface.go -destination=../../mocks/discussionContainerService_mock.go

type DiscussionContainerService interface {
	// Get a discussion container from the database by its ID
	GetDiscussionContainer(id uint) (*models.DiscussionContainer, error)
}
