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
	_, _, member, err := memberService.CreateMember(&memberForm, &userTags)

	// verify that the member object was created correctly
	assert.Equal(t, exampleMember.FirstName, member.FirstName)
	assert.Equal(t, exampleMember.LastName, member.LastName)
	assert.Equal(t, exampleMember.Email, member.Email)
	assert.Equal(t, exampleMember.Institution, member.Institution)

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
	_, _, _, err := memberService.CreateMember(&memberForm, &userTags)

	// verify the error was returned correctly
	assert.Equal(t, expectedErr, err)
}

func TestCreateMemberDuplicateEmail(t *testing.T) {
	beforeEachMember(t)

	expectedErr := fmt.Errorf("error")

	// set up repository mock to return an error
	mockMemberRepository.EXPECT().Create(gomock.Any()).Return(expectedErr)
	mockMemberRepository.EXPECT().Query(&models.Member{Email: "john.smith@gmail.com"}).Return([]*models.Member{nil, nil}, nil)

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
	_, _, _, err := memberService.CreateMember(&memberForm, &userTags)

	// verify the error was returned correctly
	assert.NotNil(t, err)
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

func TestRefreshTokenSuccess(t *testing.T) {
	beforeEachMember(t)

	// generate valid refresh token
	_, refreshToken, err := memberService.generateTokenPair(uint(1))
	assert.Nil(t, err)

	mockMemberRepository.EXPECT().GetByID(uint(1)).Return(nil, nil)

	validForm := forms.TokenRefreshForm{RefreshToken: refreshToken}

	tokenPair, err := memberService.RefreshToken(&validForm)
	assert.Nil(t, err)
	assert.NotNil(t, tokenPair.AccessToken)
	assert.NotNil(t, tokenPair.RefreshToken)
}

func TestRefreshTokenFailureWrongTyp(t *testing.T) {
	beforeEachMember(t)

	// generate valid refresh token
	accessToken, _, err := memberService.generateTokenPair(uint(1))
	assert.Nil(t, err)

	mockMemberRepository.EXPECT().GetByID(uint(1))

	validForm := forms.TokenRefreshForm{RefreshToken: accessToken}

	_, err = memberService.RefreshToken(&validForm)
	assert.NotNil(t, err)
}

func TestRefreshTokenFailureInvalid(t *testing.T) {
	beforeEachMember(t)

	mockMemberRepository.EXPECT().GetByID(uint(1))

	validForm := forms.TokenRefreshForm{RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}

	_, err := memberService.RefreshToken(&validForm)
	assert.NotNil(t, err)
}

func TestRefreshTokenFailureCantParse(t *testing.T) {
	beforeEachMember(t)

	mockMemberRepository.EXPECT().GetByID(uint(1))

	validForm := forms.TokenRefreshForm{RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6ddpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWFsIjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}

	_, err := memberService.RefreshToken(&validForm)
	assert.NotNil(t, err)
}
