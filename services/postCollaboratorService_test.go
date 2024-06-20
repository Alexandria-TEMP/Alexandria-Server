package services

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	mockPostCollaboratorRepository = mocks.NewMockModelRepositoryInterface[*models.PostCollaborator](mockCtrl)
	mockPostRepository = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)

	// Setup SUT
	postCollaboratorService = PostCollaboratorService{
		PostCollaboratorRepository: mockPostCollaboratorRepository,
		MemberRepository:           mockMemberRepository,
		PostRepository:             mockPostRepository,
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
	mockMemberRepository.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).Times(1)
	mockMemberRepository.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).Times(1)

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
	mockPostCollaboratorRepository.EXPECT().GetByID(uint(10)).Return(&models.PostCollaborator{
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

func TestMergeContributors(t *testing.T) {
	postCollaboratorServiceSetup(t)
	t.Cleanup(postCollaboratorServiceTeardown)

	postCollaborator1 := &models.PostCollaborator{
		MemberID:          1,
		Member:            models.Member{Model: gorm.Model{ID: 1}},
		CollaborationType: models.Contributor,
	}
	postCollaborator2 := &models.PostCollaborator{
		MemberID:          2,
		Member:            models.Member{Model: gorm.Model{ID: 2}},
		CollaborationType: models.Reviewer,
	}
	postCollaborator2again := &models.PostCollaborator{
		Member:            models.Member{Model: gorm.Model{ID: 2}},
		CollaborationType: models.Contributor,
	}
	branchCollaborator1 := &models.BranchCollaborator{
		MemberID: 1,
	}
	branchCollaborator2 := &models.BranchCollaborator{
		MemberID: 2,
	}
	projectPostBefore := &models.ProjectPost{
		Post: models.Post{
			Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2},
		},
	}
	projectPostAfter := &models.ProjectPost{
		Post: models.Post{
			Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2, postCollaborator2again},
		},
	}

	mockMemberRepository.EXPECT().GetByID(uint(1)).Return(&models.Member{Model: gorm.Model{ID: 1}}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(2)).Return(&models.Member{Model: gorm.Model{ID: 2}}, nil)

	mockPostRepository.EXPECT().GetByID(gomock.Any()).Return(&models.Post{
		Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2},
	}, nil).Times(1)
	mockPostRepository.EXPECT().Update(&models.Post{
		Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2, postCollaborator2again},
	}).Return(nil, nil).Times(1)

	assert.Nil(t, postCollaboratorService.MergeContributors(projectPostBefore, []*models.BranchCollaborator{branchCollaborator1, branchCollaborator2}))
	assert.Equal(t, projectPostAfter, projectPostBefore)
}

func TestMergeReviewers(t *testing.T) {
	postCollaboratorServiceSetup(t)
	t.Cleanup(postCollaboratorServiceTeardown)

	postCollaborator1 := &models.PostCollaborator{
		MemberID:          1,
		Member:            models.Member{Model: gorm.Model{ID: 1}},
		CollaborationType: models.Contributor,
	}
	postCollaborator1again := &models.PostCollaborator{
		Member:            models.Member{Model: gorm.Model{ID: 1}},
		CollaborationType: models.Reviewer,
	}
	postCollaborator2 := &models.PostCollaborator{
		MemberID:          2,
		Member:            models.Member{Model: gorm.Model{ID: 2}},
		CollaborationType: models.Reviewer,
	}
	review1 := &models.BranchReview{
		MemberID: 1,
	}
	review2 := &models.BranchReview{
		MemberID: 2,
	}
	projectPostBefore := &models.ProjectPost{
		Post: models.Post{
			Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2},
		},
	}
	projectPostAfter := &models.ProjectPost{
		Post: models.Post{
			Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2, postCollaborator1again},
		},
	}

	mockPostRepository.EXPECT().GetByID(gomock.Any()).Return(&models.Post{
		Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2},
	}, nil).Times(1)
	mockPostRepository.EXPECT().Update(&models.Post{
		Collaborators: []*models.PostCollaborator{postCollaborator1, postCollaborator2, postCollaborator1again},
	}).Return(nil, nil).Times(1)

	mockMemberRepository.EXPECT().GetByID(uint(1)).Return(&models.Member{Model: gorm.Model{ID: 1}}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(2)).Return(&models.Member{Model: gorm.Model{ID: 2}}, nil)

	assert.Nil(t, postCollaboratorService.MergeReviewers(projectPostBefore, []*models.BranchReview{review1, review2}))

	assert.Equal(t, projectPostAfter, projectPostBefore)
}
