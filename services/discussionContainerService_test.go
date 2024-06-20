package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var discussionContainerService DiscussionContainerService

func discussionContainerServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockDiscussionContainerRepository = mocks.NewMockModelRepositoryInterface[*models.DiscussionContainer](mockCtrl)

	// Setup data
	memberA = models.Member{
		Model: gorm.Model{ID: 1},
	}
	discussionA = models.Discussion{
		Model:       gorm.Model{ID: 5},
		ContainerID: 10,
		Member:      &memberA,
		MemberID:    &memberA.ID,
		Replies:     []*models.Discussion{},
		Text:        "discussion",
	}
	discussionContainerA = models.DiscussionContainer{
		Model: gorm.Model{ID: 10},
		Discussions: []*models.Discussion{
			&discussionA,
		},
	}

	// Setup SUT
	discussionContainerService = DiscussionContainerService{
		DiscussionContainerRepository: mockDiscussionContainerRepository,
	}
}

func discussionContainerServiceTeardown() {

}

func TestGetDiscussionContainer(t *testing.T) {
	discussionContainerServiceSetup(t)
	t.Cleanup(discussionContainerServiceTeardown)

	mockDiscussionContainerRepository.EXPECT().GetByID(discussionContainerA.ID).Return(&discussionContainerA, nil).Times(1)

	// Function under test
	fetchedDiscussionContainer, err := discussionContainerService.GetDiscussionContainer(discussionContainerA.ID)
	assert.Nil(t, err)
	assert.Equal(t, fetchedDiscussionContainer, &discussionContainerA)
}
