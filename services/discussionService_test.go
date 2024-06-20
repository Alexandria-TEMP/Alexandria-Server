package services

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var discussionService DiscussionService

func discussionServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockDiscussionRepository = mocks.NewMockModelRepositoryInterface[*models.Discussion](mockCtrl)
	mockDiscussionContainerRepository = mocks.NewMockModelRepositoryInterface[*models.DiscussionContainer](mockCtrl)
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)

	// Setup data
	memberA = models.Member{
		Model: gorm.Model{ID: 5},
	}
	discussionA = models.Discussion{
		Model:       gorm.Model{ID: 99},
		ContainerID: 10,
		Member:      &memberA,
		MemberID:    &memberA.ID,
		Replies:     []*models.Discussion{},
		Text:        "root discussion",
	}
	discussionContainerA = models.DiscussionContainer{
		Model:       gorm.Model{ID: 10},
		Discussions: []*models.Discussion{&discussionA},
	}

	// Setup default mock function return values
	mockMemberRepository.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).AnyTimes()
	mockDiscussionContainerRepository.EXPECT().GetByID(discussionContainerA.ID).Return(&discussionContainerA, nil).AnyTimes()
	mockDiscussionRepository.EXPECT().GetByID(discussionA.ID).Return(&discussionA, nil).AnyTimes()

	// Setup SUT
	discussionService = DiscussionService{
		DiscussionRepository:          mockDiscussionRepository,
		DiscussionContainerRepository: mockDiscussionContainerRepository,
		MemberRepository:              mockMemberRepository,
	}
}

func discussionServiceTeardown() {

}

func TestGetDiscussion(t *testing.T) {
	discussionServiceSetup(t)
	t.Cleanup(discussionServiceTeardown)

	mockDiscussionRepository.EXPECT().GetByID(uint(20)).Return(&models.Discussion{
		Model:       gorm.Model{ID: 5},
		ContainerID: 10,
		Replies:     []*models.Discussion{},
		Text:        "My Cool Discussion",
	}, nil).Times(1)

	// Function under test
	fetchedDiscussion, err := discussionService.GetDiscussion(20)
	if err != nil {
		t.Fatal(err)
	}

	expectedDiscussion := &models.Discussion{
		Model:       gorm.Model{ID: 5},
		ContainerID: 10,
		Replies:     []*models.Discussion{},
		Text:        "My Cool Discussion",
	}

	if !reflect.DeepEqual(fetchedDiscussion, expectedDiscussion) {
		t.Fatalf("fetched discussion\n%+v\nshould have equaled expected discussion\n%+v", fetchedDiscussion, expectedDiscussion)
	}
}

func TestCreateRootDiscussionGoodWeather(t *testing.T) {
	discussionServiceSetup(t)
	t.Cleanup(discussionServiceTeardown)

	discussionCreationForm := forms.RootDiscussionCreationForm{
		ContainerID: discussionContainerA.ID,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "lorem ipsum",
		},
	}

	// Setup mock function return values
	mockDiscussionRepository.EXPECT().Create(&models.Discussion{
		ContainerID: discussionContainerA.ID,
		Member:      &memberA,
		Replies:     []*models.Discussion{},
		Text:        "lorem ipsum",
	}).Return(nil).Times(1)

	// Function under test
	createdDiscussion, err := discussionService.CreateRootDiscussion(&discussionCreationForm, &memberA)
	if err != nil {
		t.Fatal(err)
	}

	expectedDiscussion := &models.Discussion{
		ContainerID: discussionContainerA.ID,
		Member:      &memberA,
		Replies:     []*models.Discussion{},
		Text:        "lorem ipsum",
	}

	if !reflect.DeepEqual(createdDiscussion, expectedDiscussion) {
		t.Fatalf("created discussion\n%+v\nshould have equaled expected discussion\n%+v", createdDiscussion, expectedDiscussion)
	}
}

func TestCreateRootDiscussionContainerDNE(t *testing.T) {
	discussionServiceSetup(t)
	t.Cleanup(discussionServiceTeardown)

	// Use a different ID that doesn't have a mock return value set up!
	containerID := discussionContainerA.ID + 5

	discussionCreationForm := forms.RootDiscussionCreationForm{
		ContainerID: containerID,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "lorem ipsum",
		},
	}

	mockDiscussionContainerRepository.EXPECT().GetByID(containerID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := discussionService.CreateRootDiscussion(&discussionCreationForm, &memberA)
	if err == nil {
		t.Fatal("creating discussion should have returned error")
	}
}

func TestCreateReplyDiscussionGoodWeather(t *testing.T) {
	discussionServiceSetup(t)
	t.Cleanup(discussionServiceTeardown)

	replyDiscussionCreationForm := forms.ReplyDiscussionCreationForm{
		ParentID: discussionA.ID,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "reply discussion",
		},
	}

	// Setup mock function return values
	mockDiscussionRepository.EXPECT().Create(&models.Discussion{
		ContainerID: discussionContainerA.ID,
		Member:      &memberA,
		Replies:     []*models.Discussion{},
		ParentID:    &discussionA.ID,
		Text:        "reply discussion",
	}).Return(nil).Times(1)

	// Function under test
	createdDiscussion, err := discussionService.CreateReply(&replyDiscussionCreationForm, &memberA)
	if err != nil {
		t.Fatal(err)
	}

	expectedDiscussion := &models.Discussion{
		ContainerID: discussionContainerA.ID,
		Member:      &memberA,
		Replies:     []*models.Discussion{},
		ParentID:    &discussionA.ID,
		Text:        "reply discussion",
	}

	if !reflect.DeepEqual(createdDiscussion, expectedDiscussion) {
		t.Fatalf("created discussion\n%+v\nshould have equaled expected discussion\n%+v", createdDiscussion, expectedDiscussion)
	}
}

func TestCreateReplyDiscussionParentDNE(t *testing.T) {
	discussionServiceSetup(t)
	t.Cleanup(discussionServiceTeardown)

	parentID := discussionA.ID + 5

	replyDiscussionCreationForm := forms.ReplyDiscussionCreationForm{
		ParentID: parentID,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "wahhh",
		},
	}

	mockDiscussionRepository.EXPECT().GetByID(parentID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := discussionService.CreateReply(&replyDiscussionCreationForm, &memberA)
	if err == nil {
		t.Fatal("creating reply should have returned err")
	}
}
