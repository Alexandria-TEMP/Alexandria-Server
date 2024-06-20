package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func setupProjectPostController(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockProjectPostService = mocks.NewMockProjectPostService(mockCtrl)
	mockDiscussionContainerService = mocks.NewMockDiscussionContainerService(mockCtrl)
	mockPostService = mocks.NewMockPostService(mockCtrl)
	mockRenderService = mocks.NewMockRenderService(mockCtrl)

	// Setup SUT
	projectPostController = ProjectPostController{
		ProjectPostService:         mockProjectPostService,
		DiscussionContainerService: mockDiscussionContainerService,
		PostService:                mockPostService,
		RenderService:              mockRenderService,
	}

	// Setup HTTP
	responseRecorder = httptest.NewRecorder()
}

func teardownProjectPostController() {

}

func TestGetProjectPostGoodWeather(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := uint(5)
	postID := uint(10)
	createdAt := time.Now().UTC()
	updatedAt := time.Now().Add(time.Hour).UTC()

	projectPost := &models.ProjectPost{
		Model:                     gorm.Model{ID: projectPostID, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Post:                      models.Post{Model: gorm.Model{ID: postID}},
		PostID:                    postID,
		OpenBranches:              []*models.Branch{},
		ClosedBranches:            []*models.ClosedBranch{},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.FormalFeedback,
		PostReviewStatus:          models.Reviewed,
	}

	// Setup mocks
	mockProjectPostService.EXPECT().GetProjectPost(projectPostID).Return(projectPost, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%d", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseDTO := &models.ProjectPostDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedDTO := &models.ProjectPostDTO{
		ID:                        projectPostID,
		PostID:                    postID,
		OpenBranchIDs:             []uint{},
		ClosedBranchIDs:           []uint{},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.FormalFeedback,
		PostReviewStatus:          models.Reviewed,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}

	assert.Equal(t, expectedDTO, responseDTO)
}

func TestGetProjectPostMalformedID(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := "oh no!!"

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%s", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostNotFound(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := uint(5)

	// Setup mocks
	mockProjectPostService.EXPECT().GetProjectPost(projectPostID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%d", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateProjectPostGoodWeather(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	form := &forms.ProjectPostCreationForm{
		AuthorMemberIDs:           []uint{0},
		Title:                     "my cool project post",
		Anonymous:                 true,
		ScientificFieldTagIDs:     []uint{1},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.DiscussionFeedback,
	}
	body, _ := json.Marshal(form)
	member := &models.Member{}

	projectPostID := uint(10)
	postID := uint(5)

	createdAt := time.Now().UTC()
	updatedAt := time.Now().Add(time.Hour).UTC()

	// The project post that is created, when the above form is sent to post service
	projectPost := &models.ProjectPost{
		Model:                     gorm.Model{ID: projectPostID, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Post:                      models.Post{Model: gorm.Model{ID: postID}},
		PostID:                    postID,
		OpenBranches:              []*models.Branch{},
		ClosedBranches:            []*models.ClosedBranch{},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.DiscussionFeedback,
		PostReviewStatus:          models.Open,
	}

	// Setup mocks
	mockProjectPostService.EXPECT().CreateProjectPost(form, member).Return(projectPost, nil, nil)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	projectPostController.CreateProjectPost(c)

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseBody := &models.ProjectPostDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseBody); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedBody := &models.ProjectPostDTO{
		ID:                        projectPostID,
		PostID:                    postID,
		OpenBranchIDs:             []uint{},
		ClosedBranchIDs:           []uint{},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.DiscussionFeedback,
		PostReviewStatus:          models.Open,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}

	assert.Equal(t, expectedBody, responseBody)
}

func TestCreateProjectPostInvalidForm(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	form := &forms.ProjectPostCreationForm{
		AuthorMemberIDs:           []uint{},
		Title:                     "",
		Anonymous:                 false,
		ScientificFieldTagIDs:     []uint{},
		ProjectCompletionStatus:   "",
		ProjectFeedbackPreference: "",
	}

	// Marshal form
	body, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/project-posts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateProjectPostSomethingNotFound(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	form := &forms.ProjectPostCreationForm{
		AuthorMemberIDs:           []uint{1},
		Title:                     "My Project Post",
		Anonymous:                 false,
		ScientificFieldTagIDs:     []uint{1},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.DiscussionFeedback,
	}
	body, _ := json.Marshal(form)
	member := &models.Member{Model: gorm.Model{ID: 1}}

	// Setup mocks
	mockProjectPostService.EXPECT().CreateProjectPost(form, member).Return(nil, fmt.Errorf("oh no"), nil).Times(1)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	projectPostController.CreateProjectPost(c)
	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateProjectPostDatabaseFailure(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	form := &forms.ProjectPostCreationForm{
		AuthorMemberIDs:           []uint{1},
		Title:                     "My Project Post",
		Anonymous:                 false,
		ScientificFieldTagIDs:     []uint{1},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.DiscussionFeedback,
	}
	body, _ := json.Marshal(form)
	member := &models.Member{Model: gorm.Model{ID: 1}}

	// Setup mocks
	mockProjectPostService.EXPECT().CreateProjectPost(form, member).Return(nil, nil, fmt.Errorf("oh no")).Times(1)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	projectPostController.CreateProjectPost(c)

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostDiscussionContainersGoodWeather(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := uint(5)

	history := &models.DiscussionContainerProjectHistoryDTO{
		CurrentDiscussionContainerID: 5,
		MergedBranchDiscussionContainers: []models.DiscussionContainerWithBranchDTO{
			{
				DiscussionContainerID: 10,
				ClosedBranchID:        7,
			},
		},
	}

	// Setup mocks
	mockProjectPostService.EXPECT().GetDiscussionContainersFromMergeHistory(projectPostID).Return(history, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%d/all-discussion-containers", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseBody := &models.DiscussionContainerProjectHistoryDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseBody); err != nil {
		t.Fatal(err)
	}

	// Check body
	assert.Equal(t, history, responseBody)
}

func TestGetProjectPostDiscussionContainersPostDNE(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := uint(10)

	// Setup mocks
	mockProjectPostService.EXPECT().GetDiscussionContainersFromMergeHistory(projectPostID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%d/all-discussion-containers", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostDiscussionContainersBadID(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := "bad iddd"

	// Setup mocks
	mockProjectPostService.EXPECT().GetDiscussionContainersFromMergeHistory(projectPostID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%s/all-discussion-containers", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostBranchesByStatusGoodWeather(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := uint(5)

	branches := &models.BranchesGroupedByReviewStatusDTO{
		OpenBranchIDs:           []uint{1},
		RejectedClosedBranchIDs: []uint{},
		ApprovedClosedBranchIDs: []uint{5, 7},
	}

	// Setup mocks
	mockProjectPostService.EXPECT().GetBranchesGroupedByReviewStatus(projectPostID).Return(branches, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%d/branches-by-status", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseBody := &models.BranchesGroupedByReviewStatusDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseBody); err != nil {
		t.Fatal(err)
	}

	// Check body
	assert.Equal(t, branches, responseBody)
}

func TestGetProjectPostBranchesByStatusPostDNE(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := uint(10)

	// Setup mocks
	mockProjectPostService.EXPECT().GetBranchesGroupedByReviewStatus(projectPostID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%d/branches-by-status", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostBranchesByStatusBadID(t *testing.T) {
	setupProjectPostController(t)
	t.Cleanup(teardownProjectPostController)

	// Setup data
	projectPostID := "bad id!!!!!!"

	// Setup mocks
	mockProjectPostService.EXPECT().GetBranchesGroupedByReviewStatus(projectPostID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/project-posts/%s/branches-by-status", projectPostID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}
