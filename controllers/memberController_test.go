package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	memberController = &MemberController{MemberService: mockMemberService, TagService: mockTagService}
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

	mockMemberService.EXPECT().CreateMember(&exampleMemberForm, gomock.Any()).Return(&exampleMember, nil).Times(1)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil).Times(1)

	exampleMemberFormJSON, _ := json.Marshal(exampleMemberForm)
	req, _ := http.NewRequest("POST", "/api/v2/members", bytes.NewBuffer(exampleMemberFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.MemberDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleMemberDTO, responsemember)
}

func TestCreateMember400(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().CreateMember(gomock.Any(), gomock.Any()).Return(&exampleMember, errors.New("some error")).Times(0)

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

func TestUpdateMember200(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().UpdateMember(&exampleMemberDTO, gomock.Any()).Return(nil).Times(1)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil).Times(1)

	exampleMemberDTOJSON, _ := json.Marshal(exampleMemberDTO)
	req, _ := http.NewRequest("PUT", "/api/v2/members", bytes.NewBuffer(exampleMemberDTOJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUpdateMember400(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().UpdateMember(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	badMemberFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("PUT", "/api/v2/members", bytes.NewBuffer(badMemberFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}
func TestUpdateMember404(t *testing.T) {
	beforeEachMember(t)

	mockMemberService.EXPECT().UpdateMember(gomock.Any(), gomock.Any()).Return(errors.New("some error")).Times(1)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil).Times(1)

	exampleMemberDTOJSON, _ := json.Marshal(exampleMemberDTO)
	req, _ := http.NewRequest("PUT", "/api/v2/members", bytes.NewBuffer(exampleMemberDTOJSON))
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
