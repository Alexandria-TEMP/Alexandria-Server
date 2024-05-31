package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	gomock "go.uber.org/mock/gomock"
)

func beforeEachMember(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockTagService = mocks.NewMockTagService(mockCtrl)
	mockMemberRepository = mocks.NewMockRepositoryInterface[*models.Member](mockCtrl)

	//need to mock repository here, how?
	memberService = MemberService {
		MemberRepository: mockMemberRepository,
	}
}


func testGetMemberSuccessful(t *testing.T, id uint) {
	id = 5

	// mock member repository here to return member when get by id called
	mockMemberRepository.EXPECT().GetByID(1).Return(exampleMember, nil)

	// call service method
	member, err := memberService.GetMember(id)
	// assert member was returned correctly
	assert.Equal(t, exampleMember, member)
	// assert there was no error
	assert.Nil(t, err)
}

func testGetMemberUnsuccessful(t *testing.T, id uint) {
	id = 5
	expectedErr := fmt.Errorf("error")
	// mock member repository here to return member when get by id called
	mockMemberRepository.EXPECT().GetByID(id).Return(nil, expectedErr)

	// call service method
	member, err := memberService.GetMember(id)
	// assert member was returned correctly
	assert.NotEqual(t, exampleMember, member)
	// assert there was no error
	assert.Equal(t, expectedErr, err)
}

func TestCreateUserSuccessful(t *testing.T) {    
	// set up repository mock to create members correctly
	mockMemberRepository.EXPECT().Create(exampleMember).Return(exampleMember, nil)

	// set up a member creation form
	memberForm := forms.MemberCreationForm {
		FirstName:		"John",
		LastName:		"Smith",
		Email:			"john.smith@gmail.com",
		Password:		"password",
		Institution:	"TU Delft",
	}

	// manually set up the member tags
	tags := []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2}

	// call service method under test
	member, err := memberService.CreateMember(&memberForm, tags)

    // verify that the member object was created correctly
	assert.Equal(t, exampleMember, member)
	// verify that there was no error
	assert.Nil(t, err)
}

func TestCreateUserUnsuccessful(t *testing.T) {
    expectedErr := fmt.Errorf("error")

	// set up repository mock to create members correctly
	mockMemberRepository.EXPECT().Create(exampleMember).Return(nil, expectedErr)

	// set up a member creation form
	memberForm := forms.MemberCreationForm {
		FirstName:		"John",
		LastName:		"Smith",
		Email:			"john.smith@gmail.com",
		Password:		"password",
		Institution:	"TU Delft",
	}

	// manually set up the member tags
	tags := []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2}

	// call service method under test
	member, err := memberService.CreateMember(&memberForm, tags)

    // verify that the member object was not created
	assert.Nil(t, member)
	// verify the error was returned correctly
	assert.Equal(t, expectedErr, err)
}