package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func beforeEachMember(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)

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
	userTags := models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{exampleSTag1, exampleSTag2},
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
	userTags := models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{exampleSTag1, exampleSTag2},
	}

	// call service method under test
	member, err := memberService.CreateMember(&memberForm, &userTags)

	// verify that the member object was not created
	assert.Nil(t, member)
	// verify the error was returned correctly
	assert.Equal(t, expectedErr, err)
}

func TestDeleteMemberSuccessful(t *testing.T) {
	beforeEachMember(t)

	id := 5

	// mock member repository here to return member when get by id called
	mockMemberRepository.EXPECT().Delete(uint(id)).Return(nil)

	// call service method
	err := memberService.DeleteMember(uint(id))

	// assert there was no error
	assert.Nil(t, err)
}

func TestDeleteMemberUnsuccessful(t *testing.T) {
	beforeEachMember(t)
	// var id uint
	id := 5
	expectedErr := fmt.Errorf("error")
	// mock member repository here to return member when get by id called
	mockMemberRepository.EXPECT().Delete(uint(id)).Return(expectedErr)

	// call service method
	err := memberService.DeleteMember(uint(id))

	// assert there expected error was returned
	assert.Equal(t, expectedErr, err)
}

func TestGetAllMembersGoodWeather(t *testing.T) {
	beforeEachMember(t)

	// Setup data
	members := []*models.Member{
		{
			Model:     gorm.Model{ID: 5},
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Model:     gorm.Model{ID: 10},
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	// Setup mocks
	mockMemberRepository.EXPECT().Query(gomock.Any()).Return(members, nil).Times(1)

	// Function under test
	actualShortFormMemberDTOs, err := memberService.GetAllMembers()
	if err != nil {
		t.Fatal(err)
	}

	expectedShortFormMemberDTOs := []*models.MemberShortFormDTO{
		{
			ID:        5,
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			ID:        10,
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	assert.Equal(t, expectedShortFormMemberDTOs, actualShortFormMemberDTOs)
}
