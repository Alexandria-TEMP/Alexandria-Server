package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func beforeEachBranch(t *testing.T) {
	t.Helper()

	_ = lock.Lock()
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
	mockBranchCollaboratorService = mocks.NewMockBranchCollaboratorService(mockCtrl)
	branchController = BranchController{
		BranchService:             mockBranchService,
		RenderService:             mockRenderService,
		BranchCollaboratorService: mockBranchCollaboratorService,
	}

	responseRecorder = httptest.NewRecorder()
}

func TestGetBranch200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetBranch(uint(1)).Return(&exampleBranch, nil)

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

	mockBranchService.EXPECT().GetBranch(uint(1)).Return(&exampleBranch, errors.New("branch not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateBranch200(t *testing.T) {
	beforeEachBranch(t)

	updatedPostTitle := "new post title"
	updatedCompletionStatus := models.Completed
	updatedFeedbackPreferences := models.DiscussionFeedback
	form := forms.BranchCreationForm{
		UpdatedPostTitle:           &updatedPostTitle,
		UpdatedCompletionStatus:    &updatedCompletionStatus,
		UpdatedFeedbackPreferences: &updatedFeedbackPreferences,
		UpdatedScientificFieldIDs:  []uint{},
		CollaboratingMemberIDs:     []uint{1},
		ProjectPostID:              5,
		BranchTitle:                "test",
	}
	body, _ := json.Marshal(form)
	member := &models.Member{}

	mockBranchService.EXPECT().CreateBranch(&form, member).Return(&exampleBranch, nil, nil)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateBranch(c)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCreateBranch4001(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{}

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)

	branchController.CreateBranch(c)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateBranch4002(t *testing.T) {
	beforeEachBranch(t)

	updatedPostTitle := "new post title"
	updatedCompletionStatus := models.Completed
	updatedFeedbackPreferences := models.DiscussionFeedback
	form := forms.BranchCreationForm{
		UpdatedPostTitle:           &updatedPostTitle,
		UpdatedCompletionStatus:    &updatedCompletionStatus,
		UpdatedFeedbackPreferences: &updatedFeedbackPreferences,
		UpdatedScientificFieldIDs:  []uint{},
		CollaboratingMemberIDs:     []uint{1},
		ProjectPostID:              5,
		BranchTitle:                "test",
	}
	body, _ := json.Marshal(form)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateBranch(c)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateBranch404(t *testing.T) {
	beforeEachBranch(t)

	updatedPostTitle := "post title"
	updatedCompletionStatus := models.Completed
	updatedFeedbackPreferences := models.DiscussionFeedback
	form := forms.BranchCreationForm{
		UpdatedPostTitle:           &updatedPostTitle,
		UpdatedCompletionStatus:    &updatedCompletionStatus,
		UpdatedFeedbackPreferences: &updatedFeedbackPreferences,
		UpdatedScientificFieldIDs:  []uint{},
		CollaboratingMemberIDs:     []uint{1},
		ProjectPostID:              5,
		BranchTitle:                "test",
	}
	body, _ := json.Marshal(form)
	member := &models.Member{}

	mockBranchService.EXPECT().CreateBranch(&form, member).Return(&exampleBranch, errors.New("failed"), nil)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateBranch(c)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateBranch500(t *testing.T) {
	beforeEachBranch(t)

	updatedPostTitle := "post title"
	updatedCompletionStatus := models.Completed
	updatedFeedbackPreferences := models.DiscussionFeedback
	form := forms.BranchCreationForm{
		UpdatedPostTitle:           &updatedPostTitle,
		UpdatedCompletionStatus:    &updatedCompletionStatus,
		UpdatedFeedbackPreferences: &updatedFeedbackPreferences,
		UpdatedScientificFieldIDs:  []uint{},
		CollaboratingMemberIDs:     []uint{1},
		ProjectPostID:              5,
		BranchTitle:                "test",
	}
	body, _ := json.Marshal(form)
	member := &models.Member{}

	mockBranchService.EXPECT().CreateBranch(&form, member).Return(&exampleBranch, nil, errors.New("failed"))

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateBranch(c)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
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

func TestGetReview200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetReview(uint(1)).Return(&exampleReview, nil)

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

	mockBranchService.EXPECT().GetReview(uint(1)).Return(&exampleReview, errors.New("branchreview not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/reviews/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateReview200(t *testing.T) {
	beforeEachBranch(t)

	form := forms.ReviewCreationForm{
		BranchReviewDecision: models.Approved,
		BranchID:             1,
	}
	body, _ := json.Marshal(form)
	member := &models.Member{}

	mockBranchService.EXPECT().CreateReview(form, member).Return(&exampleReview, nil)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateReview(c)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCreateReview4001(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{}

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)

	branchController.CreateReview(c)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateReview4002(t *testing.T) {
	beforeEachBranch(t)

	form := forms.ReviewCreationForm{
		BranchReviewDecision: models.Approved,
		BranchID:             1,
	}
	body, _ := json.Marshal(form)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateReview(c)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateReview404(t *testing.T) {
	beforeEachBranch(t)

	form := forms.ReviewCreationForm{
		BranchReviewDecision: models.Approved,
		BranchID:             1,
	}
	body, _ := json.Marshal(form)
	member := &models.Member{}

	mockBranchService.EXPECT().CreateReview(form, member).Return(&exampleReview, errors.New("branch not found"))

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	branchController.CreateReview(c)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(uint(1), uint(1)).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/can-review/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview400BranchID(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(gomock.Any(), gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/can-review/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview400MemberID(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(gomock.Any(), gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/can-review/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestMemberCanReview404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().MemberCanReview(uint(1), uint(1)).Return(false, errors.New("branch or member not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/members/1/can-review", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetBranchCollaborator200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(1)).Return(&exampleCollaborator, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/collaborators/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetBranchCollaborator400(t *testing.T) {
	beforeEachBranch(t)

	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(gomock.Any()).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/branches/collaborators/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetBranchCollaborator404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(1)).Return(nil, errors.New("collaborator not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/collaborators/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetRender200(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("../utils/test_files/good_repository_setup/render/1234.html", lock, nil, nil, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "text/html", responseRecorder.Header().Get("Content-Type"))
	assert.False(t, lock.Locked())
}

func TestGetRender202(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("", nil, errors.New("pending"), nil, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusAccepted, responseRecorder.Result().StatusCode)
}

func TestGetRender204(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("", nil, nil, errors.New("failed"), nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNoContent, responseRecorder.Result().StatusCode)
}

func TestGetRender400(t *testing.T) {
	beforeEachBranch(t)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender404(t *testing.T) {
	beforeEachBranch(t)

	mockRenderService.EXPECT().GetRenderFile(uint(1)).Return("", nil, nil, nil, errors.New("render not found"))

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetProject200(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetProject(uint(1)).Return("../utils/test_files/good_repository_setup/quarto_project.zip", lock, nil)

	req, _ := http.NewRequest("GET", "/api/v2/branches/1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "application/zip", responseRecorder.Header().Get("Content-Type"))
	assert.False(t, lock.Locked())
}

func TestGetProject400(t *testing.T) {
	beforeEachBranch(t)

	req, _ := http.NewRequest("GET", "/api/v2/branches/bad/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetProject404(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().GetProject(uint(1)).Return("", nil, errors.New("project not found"))

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

	req, _ := http.NewRequest("POST", "/api/v2/branches/1/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUploadProject400NoFile(t *testing.T) {
	beforeEachBranch(t)

	mockBranchService.EXPECT().UploadProject(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	req, _ := http.NewRequest("POST", "/api/v2/branches/1/upload", http.NoBody)
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

	req, _ := http.NewRequest("POST", "/api/v2/branches/bad/upload", body)
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

	req, _ := http.NewRequest("POST", "/api/v2/branches/1/upload", body)
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

func TestGetAllBranchCollaboratorsGoodWeather(t *testing.T) {
	beforeEachBranch(t)

	branchID := uint(15)

	branch := &models.Branch{
		Model: gorm.Model{ID: branchID},
		Collaborators: []*models.BranchCollaborator{
			{
				Model:    gorm.Model{ID: 10},
				Member:   models.Member{Model: gorm.Model{ID: 56}},
				MemberID: 56,
				BranchID: branchID,
			},
			{
				Model:    gorm.Model{ID: 20},
				Member:   models.Member{Model: gorm.Model{ID: 60}},
				MemberID: 60,
				BranchID: branchID,
			},
		},
	}

	// Setup mocks
	mockBranchService.EXPECT().GetBranch(branchID).Return(branch, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/collaborators/all/%d", branchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Read body
	bytesJSON, err := io.ReadAll(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Parse body
	responseBranchCollaboratorDTOs := []*models.BranchCollaboratorDTO{}
	if err := json.Unmarshal(bytesJSON, &responseBranchCollaboratorDTOs); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedBranchCollaboratorDTOs := []*models.BranchCollaboratorDTO{
		{
			ID:       10,
			MemberID: 56,
			BranchID: branchID,
		},
		{
			ID:       20,
			MemberID: 60,
			BranchID: branchID,
		},
	}

	assert.Equal(t, expectedBranchCollaboratorDTOs, responseBranchCollaboratorDTOs)
}

func TestGetAllBranchCollaboratorsBranchDNE(t *testing.T) {
	beforeEachBranch(t)

	branchID := uint(20)

	// Setup mocks
	mockBranchService.EXPECT().GetBranch(branchID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/collaborators/all/%d", branchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetAllBranchCollaboratorsBadBranchID(t *testing.T) {
	beforeEachBranch(t)

	// Construct request
	req, err := http.NewRequest("GET", "/api/v2/branches/collaborators/all/badbranchID", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReviewStatusGoodWeather(t *testing.T) {
	beforeEachBranch(t)

	branchID := uint(10)

	// Setup mocks
	mockBranchService.EXPECT().GetAllBranchReviewStatuses(branchID).Return([]models.BranchReviewDecision{
		models.Approved,
		models.Approved,
		models.Rejected,
	}, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/%d/review-statuses", branchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Get JSON
	jsonBytes, err := io.ReadAll(responseRecorder.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	// Parse JSON
	var responseReviewStatuses []models.BranchReviewDecision
	if err := json.Unmarshal(jsonBytes, &responseReviewStatuses); err != nil {
		t.Fatal(err)
	}

	// Check JSON
	expectedReviewStatuses := []models.BranchReviewDecision{
		models.Approved,
		models.Approved,
		models.Rejected,
	}

	assert.Equal(t, expectedReviewStatuses, responseReviewStatuses)
}

func TestGetReviewStatusBranchDNE(t *testing.T) {
	beforeEachBranch(t)

	branchID := uint(10)

	// Setup mocks
	mockBranchService.EXPECT().GetAllBranchReviewStatuses(branchID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/%d/review-statuses", branchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetClosedBranch(t *testing.T) {
	beforeEachBranch(t)

	// Setup data
	branchID := uint(8)
	closedBranchID := uint(5)
	supercededBranchID := uint(10)
	projectPostID := uint(2)

	closedBranch := &models.ClosedBranch{
		Model:                gorm.Model{ID: closedBranchID},
		Branch:               models.Branch{},
		BranchID:             branchID,
		SupercededBranch:     &models.Branch{},
		SupercededBranchID:   &supercededBranchID,
		ProjectPostID:        projectPostID,
		BranchReviewDecision: models.Approved,
	}

	// Setup mocks
	mockBranchService.EXPECT().GetClosedBranch(closedBranchID).Return(closedBranch, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/closed/%d", closedBranchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseClosedBranchDTO := &models.ClosedBranchDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseClosedBranchDTO); err != nil {
		t.Fatal(err)
	}

	expectedClosedBranchDTO := &models.ClosedBranchDTO{
		ID:                   closedBranchID,
		BranchID:             branchID,
		SupercededBranchID:   &supercededBranchID,
		ProjectPostID:        projectPostID,
		BranchReviewDecision: models.Approved,
	}

	// Check body
	assert.Equal(t, expectedClosedBranchDTO, responseClosedBranchDTO)
}

func TestGetClosedBranchDNE(t *testing.T) {
	beforeEachBranch(t)

	// Setup data
	closedBranchID := uint(10)

	// Setup mocks
	mockBranchService.EXPECT().GetClosedBranch(closedBranchID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/closed/%d", closedBranchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetClosedBranchInvalidID(t *testing.T) {
	beforeEachBranch(t)

	// Setup data
	closedBranchID := "Bad!!!"

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/branches/closed/%s", closedBranchID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}
