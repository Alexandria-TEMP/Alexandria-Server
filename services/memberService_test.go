package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

func beforeEachMember(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)

	memberService = MemberService{
		Secret:           "secret",
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
	mockMemberRepository.EXPECT().Create(gomock.Any()).Return(nil)
	mockMemberRepository.EXPECT().Query(&models.Member{Email: "john.smith@gmail.com"}).Return(nil, nil)

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
	assert.Equal(t, exampleMember.FirstName, member.Member.FirstName)
	assert.Equal(t, exampleMember.LastName, member.Member.LastName)
	assert.Equal(t, exampleMember.Email, member.Member.Email)
	assert.Equal(t, exampleMember.Institution, member.Member.Institution)

	// verify that there was no error
	assert.Nil(t, err)
}

func TestCreateMemberUnsuccessful(t *testing.T) {
	beforeEachMember(t)

	expectedErr := fmt.Errorf("error")

	// set up repository mock to return an error
	mockMemberRepository.EXPECT().Create(gomock.Any()).Return(expectedErr)
	mockMemberRepository.EXPECT().Query(&models.Member{Email: "john.smith@gmail.com"}).Return(nil, nil)

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
	_, err := memberService.CreateMember(&memberForm, &userTags)

	// verify the error was returned correctly
	assert.Equal(t, expectedErr, err)
}

func TestUpdateMemberSuccessful(t *testing.T) {
	beforeEachMember(t)

	// mock member repository to return the example member
	mockMemberRepository.EXPECT().GetByID(gomock.Any()).Return(&exampleMember, nil)

	// mock member repository here to return no error
	mockMemberRepository.EXPECT().Update(&exampleMember).Return(&exampleMember, nil)

	// call service method
	err := memberService.UpdateMember(&exampleMemberDTO, &exampleMember.ScientificFieldTagContainer)
	// assert there was no error
	assert.Nil(t, err)
}

func TestUpdateMemberUnsuccessful(t *testing.T) {
	beforeEachMember(t)

	expectedErr := fmt.Errorf("error")

	// mock member repository to return the example member
	mockMemberRepository.EXPECT().GetByID(gomock.Any()).Return(&exampleMember, nil)

	// mock member repository to return an error
	mockMemberRepository.EXPECT().Update(&exampleMember).Return(&exampleMember, expectedErr)

	// call service method
	err := memberService.UpdateMember(&exampleMemberDTO, &exampleMember.ScientificFieldTagContainer)

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

func TestLoginMemberValidSuccessful(t *testing.T) {
	beforeEachMember(t)

	mockMemberRepository.EXPECT().Query(&models.Member{Email: exampleMemberAuthForm.Email}).Return([]*models.Member{&exampleMemberWithPassword}, nil)

	loggedInMember, err := memberService.LogInMember(&exampleMemberAuthForm)
	assert.Nil(t, err)
	assert.Equal(t, exampleMemberDTO, loggedInMember.Member)
}

func TestLoginMemberInvalidSuccessful(t *testing.T) {
	beforeEachMember(t)

	mockMemberRepository.EXPECT().Query(&models.Member{Email: exampleMemberAuthForm.Email}).Return([]*models.Member{&exampleMemberWithPassword}, nil)

	_, err := memberService.LogInMember(&forms.MemberAuthForm{
		Email:    "john.smith@gmail.com",
		Password: "wrong",
	})
	assert.NotNil(t, err)
}

func TestLoginMemberFailureQueryFailed(t *testing.T) {
	beforeEachMember(t)

	mockMemberRepository.EXPECT().Query(&models.Member{Email: exampleMemberAuthForm.Email}).Return(nil, errors.New("failed"))

	_, err := memberService.LogInMember(&exampleMemberAuthForm)
	assert.NotNil(t, err)
}

func TestLoginMemberFailureNoEmailMatches(t *testing.T) {
	beforeEachMember(t)

	mockMemberRepository.EXPECT().Query(&models.Member{Email: exampleMemberAuthForm.Email}).Return(nil, nil)

	_, err := memberService.LogInMember(&exampleMemberAuthForm)
	assert.NotNil(t, err)
}

// func TestRefreshTokenSuccess(t *testing.T) {

// }
