package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

func beforeEachMember(t *testing.T) {
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	responseRecorder = httptest.NewRecorder()

	mockMemberService = mock_interfaces.NewMockMemberService(mockCtrl)
	mockTagService = mock_interfaces.NewMockTagService(mockCtrl)
	memberController = MemberController{MemberService: mockMemberService, TagService: mockTagService}
}

func TestGetMember200(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(&exampleMember, nil).Times(1)

	req, _ := http.NewRequest("GET", "/api/v2/members/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.MemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleMemberDTO, responsemember)
}

func TestGetMember400(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().GetMember(gomock.Any()).Return(&exampleMember, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/members/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetMember404(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(&models.Member{}, errors.New("some error")).Times(1)

	exampleMemberJSON, _ := json.Marshal(exampleMember)
	req, _ := http.NewRequest("GET", "/api/v2/members/1", bytes.NewBuffer(exampleMemberJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateMember200(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().CreateMember(&exampleMemberForm, gomock.Any()).Return("access", int64(1), "refresh", int64(2), &exampleMember, nil).Times(1)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil).Times(1)

	exampleMemberFormJSON, _ := json.Marshal(exampleMemberForm)
	req, _ := http.NewRequest("POST", "/api/v2/members", bytes.NewBuffer(exampleMemberFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.LoggedInMemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleMemberLoggedInDTO, responsemember)
}

func TestCreateMember400(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().CreateMember(gomock.Any(), gomock.Any()).Return("", int64(0), "", int64(0), nil, errors.New("some error")).Times(0)

	badMemberFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v2/members", bytes.NewBuffer(badMemberFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestDeleteMember200(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().DeleteMember(uint(1)).Return(nil).Times(1)

	req, _ := http.NewRequest("DELETE", "/api/v2/members/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestDeleteMember400(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().DeleteMember(gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("DELETE", "/api/v2/members/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestDeleteMember404(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().DeleteMember(uint(1)).Return(errors.New("some error")).Times(1)

	req, _ := http.NewRequest("DELETE", "/api/v2/members/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetAllMembers200(t *testing.T) {
	beforeEachMember(t)
	mockMemberService.EXPECT().GetAllMembers().Return([]*models.MemberShortFormDTO{
		{ID: 3, FirstName: "eve", LastName: "eeve"},
	}, nil)

	req, _ := http.NewRequest("GET", "/api/v2/members", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	var responseMember []*models.MemberShortFormDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responseMember)

	assert.Equal(t, responseMember, []*models.MemberShortFormDTO{
		{ID: 3, FirstName: "eve", LastName: "eeve"},
	})
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetAllMembers404(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().GetAllMembers().Return(nil, errors.New("some error")).Times(1)

	req, _ := http.NewRequest("GET", "/api/v2/members", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestLoginMember200(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().LogInMember(&exampleMemberAuthForm).Return(&exampleMember, "access", int64(1), "refresh", int64(2), nil)

	exampleMemberAuthFormJSON, _ := json.Marshal(exampleMemberAuthForm)
	req, _ := http.NewRequest("POST", "/api/v2/members/login", bytes.NewBuffer(exampleMemberAuthFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.LoggedInMemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, exampleMemberLoggedInDTO, responsemember)
}

func TestLoginMember4001(t *testing.T) {
	beforeEachMember(t)

	exampleMemberAuthFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v2/members/login", bytes.NewBuffer(exampleMemberAuthFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.LoggedInMemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestLoginMember4002(t *testing.T) {
	beforeEachMember(t)

	exampleMemberAuthFormJSON, _ := json.Marshal(models.LoggedInMemberDTO{})
	req, _ := http.NewRequest("POST", "/api/v2/members/login", bytes.NewBuffer(exampleMemberAuthFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.LoggedInMemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestLoginMember401(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().LogInMember(&exampleMemberAuthForm).Return(&exampleMember, "access", int64(1), "refresh", int64(2), errors.New("failed"))

	exampleMemberAuthFormJSON, _ := json.Marshal(exampleMemberAuthForm)
	req, _ := http.NewRequest("POST", "/api/v2/members/login", bytes.NewBuffer(exampleMemberAuthFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.LoggedInMemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func TestRefreshToken200(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().RefreshToken(&exampleTokenRefreshForm).Return("5678", int64(1), "9012", int64(2), nil)

	exampleTokenRefreshFormJSON, _ := json.Marshal(exampleTokenRefreshForm)
	req, _ := http.NewRequest("POST", "/api/v2/members/token", bytes.NewBuffer(exampleTokenRefreshFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.TokenPairDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, exampleTokenPairDTO, responsemember)
}

func TestRefreshToken4001(t *testing.T) {
	beforeEachMember(t)

	exampleTokenRefreshFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v2/members/token", bytes.NewBuffer(exampleTokenRefreshFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.TokenPairDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestRefreshToken4002(t *testing.T) {
	beforeEachMember(t)

	exampleTokenRefreshFormJSON, _ := json.Marshal(forms.TokenRefreshForm{})
	req, _ := http.NewRequest("POST", "/api/v2/members/token", bytes.NewBuffer(exampleTokenRefreshFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.TokenPairDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestRefreshToken401(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().RefreshToken(&exampleTokenRefreshForm).Return("5678", int64(1), "9012", int64(2), errors.New("failed"))

	exampleTokenRefreshFormJSON, _ := json.Marshal(exampleTokenRefreshForm)
	req, _ := http.NewRequest("POST", "/api/v2/members/token", bytes.NewBuffer(exampleTokenRefreshFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.TokenPairDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func TestCreateMemberFormValidationFailed(t *testing.T) {
	beforeEachMember(t)

	form := forms.MemberCreationForm{
		FirstName:             "",
		LastName:              "",
		Email:                 "",
		Password:              "",
		Institution:           "",
		ScientificFieldTagIDs: []uint{},
	}

	// Marshal form
	body, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/members", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateMemberDatabaseFailed(t *testing.T) {
	beforeEachMember(t)

	form := forms.MemberCreationForm{
		FirstName:             "John",
		LastName:              "Doe",
		Email:                 "todo@todo.todo",
		Password:              "password",
		Institution:           "TU delft",
		ScientificFieldTagIDs: []uint{},
	}

	// Marshal form
	body, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	// Setup mocks
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil).Times(1)
	mockMemberService.EXPECT().CreateMember(&form, gomock.Any()).Return("access", int64(1), "refresh", int64(2), &exampleMember, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/members", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}
