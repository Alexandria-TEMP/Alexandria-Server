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
var branchCollaboratorService BranchCollaboratorService

func branchCollaboratorServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup data
	memberA = models.Member{
		Model: gorm.Model{ID: 10},
	}
	memberB = models.Member{
		Model: gorm.Model{ID: 15},
	}

	// Setup mocks
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	branchCollaboratorRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.BranchCollaborator](mockCtrl)

	// Setup mock function return values
	mockMemberRepository.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).AnyTimes()

	// Setup SUT
	branchCollaboratorService = BranchCollaboratorService{
		MemberRepository:             mockMemberRepository,
		BranchCollaboratorRepository: branchCollaboratorRepositoryMock,
	}
}

func branchCollaboratorServiceTeardown() {

}

func TestCreateBranchCollaboratorsGoodWeather(t *testing.T) {
	branchCollaboratorServiceSetup(t)
	t.Cleanup(branchCollaboratorServiceTeardown)

	memberIDs := []uint{memberA.ID, memberB.ID}
	anonymous := false

	// Function under test
	createdBranchCollaborators, err := branchCollaboratorService.MembersToBranchCollaborators(memberIDs, anonymous)
	if err != nil {
		t.Fatal(err)
	}

	expectedBranchCollaborators := []*models.BranchCollaborator{
		{Member: memberA},
		{Member: memberB},
	}

	if !reflect.DeepEqual(createdBranchCollaborators, expectedBranchCollaborators) {
		t.Fatalf("created branch collaborators\n%+v\ndid not equal expected branch collaborators\n%+v\n",
			createdBranchCollaborators, expectedBranchCollaborators)
	}
}

func TestCreateBranchCollaboratorsAnonymous(t *testing.T) {
	branchCollaboratorServiceSetup(t)
	t.Cleanup(branchCollaboratorServiceTeardown)

	memberIDs := []uint{memberA.ID, memberB.ID}
	anonymous := true

	// Function under test
	createdBranchCollaborators, err := branchCollaboratorService.MembersToBranchCollaborators(memberIDs, anonymous)
	if err != nil {
		t.Fatal(err)
	}

	expectedBranchCollaborators := []*models.BranchCollaborator{}

	if !reflect.DeepEqual(createdBranchCollaborators, expectedBranchCollaborators) {
		t.Fatalf("created branch collaborators\n%+v\ndid not equal expected branch collaborators\n%+v\n",
			createdBranchCollaborators, expectedBranchCollaborators)
	}
}

func TestCreateBranchCollaboratorsAtLeastOneMember(t *testing.T) {
	branchCollaboratorServiceSetup(t)
	t.Cleanup(branchCollaboratorServiceTeardown)

	memberIDs := []uint{}
	anonymous := false

	// Function under test
	_, err := branchCollaboratorService.MembersToBranchCollaborators(memberIDs, anonymous)

	if err == nil {
		t.Fatalf("creating branch collaborators with empty member list should fail")
	}
}

func TestGetBranchCollaborator(t *testing.T) {
	branchCollaboratorServiceSetup(t)
	t.Cleanup(branchCollaboratorServiceTeardown)

	// Setup mock function return values
	branchCollaboratorRepositoryMock.EXPECT().GetByID(uint(10)).Return(&models.BranchCollaborator{
		Model:  gorm.Model{ID: 5},
		Member: models.Member{Model: gorm.Model{ID: 10}},
	}, nil).Times(1)

	// Function under test
	fetchedBranchCollaborator, err := branchCollaboratorService.GetBranchCollaborator(10)
	if err != nil {
		t.Fatal(err)
	}

	expectedBranchCollaborator := &models.BranchCollaborator{
		Model:  gorm.Model{ID: 5},
		Member: models.Member{Model: gorm.Model{ID: 10}},
	}

	if !reflect.DeepEqual(fetchedBranchCollaborator, expectedBranchCollaborator) {
		t.Fatalf("fetched branch collaborator \n%+v\nshould have equaled expected branch collaborator \n%+v", fetchedBranchCollaborator, expectedBranchCollaborator)
	}
}
