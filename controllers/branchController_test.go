package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func beforeEachBranch(t *testing.T) {
	t.Helper()

	exampleBranch = models.Branch{
		Model: gorm.Model{ID: 1},
	}
	exampleReview = models.BranchReview{
		Model: gorm.Model{ID: 2},
	}
	exampleCollaborator = models.BranchCollaborator{
		Model: gorm.Model{ID: 3},
	}

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockBranchService = mocks.NewMockBranchService(mockCtrl)
	mockRenderService = mocks.NewMockRenderService(mockCtrl)
	branchController = BranchController{
		BranchService: mockBranchService,
		RenderService: mockRenderService,
	}

	responseRecorder = httptest.NewRecorder()
}

func TestGetBranch200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranch(uint(1)).Return(exampleBranch, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetBranch400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranch(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetBranch404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranch(uint(1)).Return(exampleBranch, errors.New("branch not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateBranch200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateBranch(gomock.Any()).Return(exampleBranch, nil, nil)

	form := forms.BranchCreationForm{ProjectPostID: 5}
	body, err := json.Marshal(form)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/branches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCreateBranch400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateBranch(gomock.Any()).Times(0)

	req, _ := http.NewRequest("POST", "/api/v2/branches", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateBranch404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateBranch(gomock.Any()).Return(exampleBranch, errors.New("parent branch not found"), nil)

	form := forms.BranchCreationForm{ProjectPostID: 5}
	body, err := json.Marshal(form)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/branches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateBranch500(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateBranch(gomock.Any()).Return(exampleBranch, nil, errors.New("internal server error"))

	form := forms.BranchCreationForm{ProjectPostID: 5}
	body, err := json.Marshal(form)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/branches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestUpdateBranch200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UpdateBranch(gomock.Any()).Return(exampleBranch, nil)

	dto := models.BranchDTO{ID: 1}
	body, err := json.Marshal(dto)
	assert.Nil(t, err)

	req, _ := http.NewRequest("PUT", "/api/v2/branches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUpdateBranch400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UpdateBranch(gomock.Any()).Times(0)

	req, _ := http.NewRequest("PUT", "/api/v2/branches", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateBranch404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UpdateBranch(gomock.Any()).Return(exampleBranch, errors.New("branch not found"))

	dto := models.BranchDTO{ID: 1}
	body, err := json.Marshal(dto)
	assert.Nil(t, err)

	req, _ := http.NewRequest("PUT", "/api/v2/branches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestDeleteBranch200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().DeleteBranch(uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v2/branches/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestDeleteBranch400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().DeleteBranch(gomock.Any()).Times(0)

	req, _ := http.NewRequest("DELETE", "/api/v2/branches/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestDeleteBranch404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().DeleteBranch(uint(1)).Return(errors.New("branch not found"))

	req, _ := http.NewRequest("DELETE", "/api/v2/branches/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetReviewStatus200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReviewStatus(uint(1)).Return([]models.BranchReviewDecision{models.Approved, models.Rejected}, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/branchreview-statuses", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetReviewStatus400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReviewStatus(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/branchreview-statuses", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReviewStatus404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReviewStatus(uint(1)).Return(nil, errors.New("branch not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/branchreview-statuses", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetReview200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReview(uint(1)).Return(exampleReview, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/reviews/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetReview400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReview(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/reviews/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReview404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReview(uint(1)).Return(exampleReview, errors.New("branchreview not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/reviews/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateReview200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateReview(gomock.Any()).Return(exampleReview, nil)

	form := forms.ReviewCreationForm{BranchID: 1}
	body, err := json.Marshal(form)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/branches/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCreateReview400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateReview(gomock.Any()).Times(0)

	req, _ := http.NewRequest("POST", "/api/v2/branches/reviews", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateReview404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().CreateReview(gomock.Any()).Return(exampleReview, errors.New("branch not found"))

	form := forms.ReviewCreationForm{BranchID: 1}
	body, err := json.Marshal(form)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/branches/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(uint(1), uint(1)).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/can-branchreview/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview400BranchID(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(gomock.Any(), gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/can-branchreview/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview400MemberID(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(gomock.Any(), gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/can-branchreview/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(uint(1), uint(1)).Return(false, errors.New("branch or member not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/members/1/can-branchreview", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetBranchCollaborator200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranchCollaborator(uint(1)).Return(&exampleCollaborator, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/collaborators/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetBranchCollaborator400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranchCollaborator(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/collaborators/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetBranchCollaborator404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranchCollaborator(uint(1)).Return(nil, errors.New("collaborator not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/collaborators/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetRender200(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "text/html", responseRecorder.Header().Get("Content-Type"))
}

func TestGetRender202(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("", errors.New("pending"), nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusAccepted, responseRecorder.Result().StatusCode)
}

func TestGetRender400(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender404(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("", nil, errors.New("render not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetProject200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetProject(uint(1)).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "application/zip", responseRecorder.Header().Get("Content-Type"))
}

func TestGetProject400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetProject(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetProject404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetProject(uint(1)).Return("", errors.New("project not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestUploadProject200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UploadProject(gomock.Any(), gomock.Any(), uint(1)).Return(nil)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.zip")
	_, _ = part.Write([]byte("file content"))

	writer.Close()

	req, _ := http.NewRequest("POST", "/api/v2/branches/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUploadProject400NoFile(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UploadProject(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	req, _ := http.NewRequest("POST", "/api/v2/branches/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUploadProject400InvalidID(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UploadProject(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.zip")
	_, _ = part.Write([]byte("file content"))

	writer.Close()

	req, _ := http.NewRequest("POST", "/api/v2/branches/bad", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUploadProject500(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UploadProject(gomock.Any(), gomock.Any(), uint(1)).Return(errors.New("upload error"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.zip")
	_, _ = part.Write([]byte("file content"))

	writer.Close()

	req, _ := http.NewRequest("POST", "/api/v2/branches/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetFiletree200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetFiletree(uint(1)).Return(map[string]int64{}, nil, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetFiletree400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetFiletree(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFiletree404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetFiletree(uint(1)).Return(nil, errors.New("filetree not found"), nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetFiletree500(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetFiletree(uint(1)).Return(nil, nil, errors.New("internal server error"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}
