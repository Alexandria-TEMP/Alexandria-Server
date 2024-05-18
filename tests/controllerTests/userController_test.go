package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

func beforeEachUser(t *testing.T) {
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	responseRecorder = httptest.NewRecorder()

	mockUserService = mock_interfaces.NewMockUserService(mockCtrl)
	userController = &controllers.UserController{UserService: mockUserService}
}

func TestGetMember200(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().GetMember(uint64(1)).Return(&exampleMember, nil).Times(1)

	req, _ := http.NewRequest("GET", "/api/v1/member/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.Member

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleMember, responsemember)
}

func TestGetMember400(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().GetMember(gomock.Any()).Return(&exampleMember, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/member/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetMember410(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().GetMember(uint64(1)).Return(&models.Member{}, errors.New("some error")).Times(1)

	exampleMemberJSON, _ := json.Marshal(exampleMember)
	req, _ := http.NewRequest("GET", "/api/v1/member/1", bytes.NewBuffer(exampleMemberJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}

func TestCreateMember200(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().CreateMember(&exampleMemberForm).Return(&exampleMember).Times(1)

	exampleMemberFormJSON, _ := json.Marshal(exampleMemberForm)
	req, _ := http.NewRequest("member", "/api/v1/member", bytes.NewBuffer(exampleMemberFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.Member

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleMember, responsemember)
}

func TestCreateMember400(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().CreateMember(gomock.Any()).Return(&exampleMember).Times(0)

	badMemberFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v1/member", bytes.NewBuffer(badMemberFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateMember200(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().UpdateMember(&exampleMember).Return(nil).Times(1)

	exampleMemberJSON, _ := json.Marshal(exampleMember)
	req, _ := http.NewRequest("PUT", "/api/v1/member", bytes.NewBuffer(exampleMemberJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUpdateMember400(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().UpdateMember(&exampleMember).Return(nil).Times(0)

	exampleMemberJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("PUT", "/api/v1/member", bytes.NewBuffer(exampleMemberJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateMember410(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().UpdateMember(&exampleMember).Return(errors.New("some error")).Times(1)

	exampleMemberJSON, _ := json.Marshal(exampleMember)
	req, _ := http.NewRequest("PUT", "/api/v1/member", bytes.NewBuffer(exampleMemberJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}

func TestGetCollaborator200(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().GetCollaborator(uint64(2)).Return(&exampleCollaborator, nil).Times(1)

	req, _ := http.NewRequest("GET", "/api/v1/collaborator/2", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.PostCollaborator

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleCollaborator, responsemember)
}

func TestGetCollaborator400(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().GetCollaborator(gomock.Any()).Return(&exampleCollaborator, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/collaborator/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetCollaborator410(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().GetCollaborator(uint64(1)).Return(&models.PostCollaborator{}, errors.New("some error")).Times(1)

	exampleCollaboratorJSON, _ := json.Marshal(exampleCollaborator)
	req, _ := http.NewRequest("GET", "/api/v1/collaborator/1", bytes.NewBuffer(exampleCollaboratorJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}

func TestCreateCollaborator200(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().CreateCollaborator(&exampleCollaboratorForm).Return(&exampleCollaborator).Times(1)

	exampleCollaboratorFormJSON, _ := json.Marshal(exampleCollaboratorForm)
	req, _ := http.NewRequest("member", "/api/v1/collaborator", bytes.NewBuffer(exampleCollaboratorFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsemember models.PostCollaborator

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsemember)

	assert.Equal(t, exampleCollaborator, responsemember)
}

func TestCreateCollaborator400(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().CreateCollaborator(gomock.Any()).Return(&exampleCollaborator).Times(0)

	badCollaboratorFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v1/collaborator", bytes.NewBuffer(badCollaboratorFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateCollaborator200(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().UpdateCollaborator(&exampleCollaborator).Return(nil).Times(1)

	exampleCollaboratorJSON, _ := json.Marshal(exampleCollaborator)
	req, _ := http.NewRequest("PUT", "/api/v1/collaborator", bytes.NewBuffer(exampleCollaboratorJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUpdateCollaborator400(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().UpdateCollaborator(gomock.Any()).Return(nil).Times(0)

	exampleCollaboratorJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("PUT", "/api/v1/collaborator", bytes.NewBuffer(exampleCollaboratorJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateCollaborator410(t *testing.T) {
	beforeEachUser(t)

	mockUserService.EXPECT().UpdateCollaborator(&exampleCollaborator).Return(errors.New("some error")).Times(1)

	exampleCollaboratorJSON, _ := json.Marshal(exampleCollaborator)
	req, _ := http.NewRequest("PUT", "/api/v1/collaborator", bytes.NewBuffer(exampleCollaboratorJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}
