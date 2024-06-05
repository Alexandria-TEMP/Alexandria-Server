package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	tags "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	gomock "go.uber.org/mock/gomock"
)

func beforeEachMember(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockMemberRepository = mocks.NewMockRepositoryInterface[*models.Member](mockCtrl)

	// need to mock repository here, how?
	memberService = MemberService{
		MemberRepository: mockMemberRepository,
	}
}

func TestGetMemberSuccessful(t *testing.T) {
	beforeEachMember(t)

	id := 5

	// mock member repository here to return member when get by id called
	mockMemberRepository.EXPECT().GetByID(uint(id)).Return(&exampleMember, nil)

	// call service method
	member, err := memberService.GetMember(uint(id))
	// assert member was returned correctly
	assert.Equal(t, &exampleMember, member)
	// assert there was no error
	assert.Nil(t, err)
}

func TestGetMemberUnsuccessful(t *testing.T) {
	beforeEachMember(t)
	// var id uint
	id := 5
	expectedErr := fmt.Errorf("error")
	// mock member repository here to return member when get by id called
	mockMemberRepository.EXPECT().GetByID(uint(id)).Return(nil, expectedErr)

	// call service method
	member, err := memberService.GetMember(uint(id))
	// assert member was returned correctly
	assert.NotEqual(t, &exampleMember, member)
	// assert there was no error
	assert.Equal(t, expectedErr, err)
}

func TestCreateMemberSuccessful(t *testing.T) {
	beforeEachMember(t)
	// set up repository mock to create members correctly
	mockMemberRepository.EXPECT().Create(&exampleMember).Return(nil)

	// set up a member creation form
	memberForm := forms.MemberCreationForm{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
	}

	// manually set up the member tags
	userTags := tags.ScientificFieldTagContainer{
		ScientificFieldTags: []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2},
	}
	// call service method under test
	member, err := memberService.CreateMember(&memberForm, &userTags)

	// verify that the member object was created correctly
	assert.Equal(t, &exampleMember, member)
	// verify that there was no error
	assert.Nil(t, err)
}

func TestCreateMemberUnsuccessful(t *testing.T) {
	beforeEachMember(t)

	expectedErr := fmt.Errorf("error")

	// set up repository mock to return an error
	mockMemberRepository.EXPECT().Create(&exampleMember).Return(expectedErr)

	// set up a member creation form
	memberForm := forms.MemberCreationForm{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
	}

	// manually set up the member tags
	userTags := tags.ScientificFieldTagContainer{
		ScientificFieldTags: []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2},
	}

	// call service method under test
	member, err := memberService.CreateMember(&memberForm, &userTags)

	// verify that the member object was not created
	assert.Nil(t, member)
	// verify the error was returned correctly
	assert.Equal(t, expectedErr, err)
}
