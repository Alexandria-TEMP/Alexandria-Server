package services

import (
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var postCollaboratorService PostCollaboratorService

func postCollaboratorServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup data
	memberA = models.Member{
		Model: gorm.Model{ID: 20},
	}
	memberB = models.Member{
		Model: gorm.Model{ID: 50},
	}

	// Setup mocks
	memberRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	postCollaboratorRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.PostCollaborator](mockCtrl)

	// Setup SUT
	postCollaboratorService = PostCollaboratorService{
		MemberRepository:           memberRepositoryMock,
		PostCollaboratorRepository: postCollaboratorRepositoryMock,
	}
}

func postCollaboratorServiceTeardown() {

}

func TestMembersToPostCollaboratorsGoodWeather(t *testing.T) {
	postCollaboratorServiceSetup(t)
	t.Cleanup(postCollaboratorServiceTeardown)

	memberIDs := []uint{memberA.ID, memberB.ID}
	anonymous := false
	collaborationType := models.Contributor

	// Setup mock function return values
	memberRepositoryMock.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).Times(1)
	memberRepositoryMock.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).Times(1)

	// Function under test
	createdPostCollaborators, err := postCollaboratorService.MembersToPostCollaborators(memberIDs, anonymous, collaborationType)

	if err != nil {
		t.Fatal(err)
	}

	expectedPostCollaborators := []*models.PostCollaborator{
		{
			Member:            memberA,
			CollaborationType: collaborationType,
		},
		{
			Member:            memberB,
			CollaborationType: collaborationType,
		},
	}

	if !reflect.DeepEqual(createdPostCollaborators, expectedPostCollaborators) {
		t.Fatalf("created post collaborators:\n%+v\ndid not equal expected post collaborators:\n%+v\n",
			createdPostCollaborators, expectedPostCollaborators)
	}
}

// Even when author IDs are passed, if it's meant to be anonymous, an empty collaborator list will be returned
func TestMembersToPostCollaboratorsAnonymous(t *testing.T) {
	postCollaboratorServiceSetup(t)
	t.Cleanup(postCollaboratorServiceTeardown)

	memberIDs := []uint{10, 50, 20, 30}
	anonymous := true
	collaborationType := models.Author

	// Function under test
	createdPostCollaborators, err := postCollaboratorService.MembersToPostCollaborators(memberIDs, anonymous, collaborationType)

	if err != nil {
		t.Fatal(err)
	}

	if len(createdPostCollaborators) > 0 {
		t.Fatalf("should not have created post collaborators, but received: %+v", createdPostCollaborators)
	}
}

// If the list is not anonymous, at least one author must be passed!
func TestMembersToPostCollaboratorsAtLeastOneAuthor(t *testing.T) {
	postCollaboratorServiceSetup(t)
	t.Cleanup(postCollaboratorServiceTeardown)

	// Function under test
	_, err := postCollaboratorService.MembersToPostCollaborators([]uint{}, false, models.Author)

	if err == nil {
		t.Fatalf("should have errored on empty author list when not anonymous")
	}
}

func TestGetPostCollaborator(t *testing.T) {
	postCollaboratorServiceSetup(t)
	t.Cleanup(postCollaboratorServiceTeardown)

	// Setup mock function return values
	postCollaboratorRepositoryMock.EXPECT().GetByID(uint(10)).Return(&models.PostCollaborator{
		Model:             gorm.Model{ID: 5},
		Member:            models.Member{Model: gorm.Model{ID: 10}},
		CollaborationType: models.Author,
	}, nil).Times(1)

	// Function under test
	fetchedPostCollaborator, err := postCollaboratorService.GetPostCollaborator(10)
	if err != nil {
		t.Fatal(err)
	}

	expectedPostCollaborator := &models.PostCollaborator{
		Model:             gorm.Model{ID: 5},
		Member:            models.Member{Model: gorm.Model{ID: 10}},
		CollaborationType: models.Author,
	}

	if !reflect.DeepEqual(fetchedPostCollaborator, expectedPostCollaborator) {
		t.Fatalf("fetched post collaborator \n%+v\nshould have equaled expected post collaborator \n%+v", fetchedPostCollaborator, expectedPostCollaborator)
	}
}
