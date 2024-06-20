package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func setupPostController(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	_ = lock.Lock()
	mockPostService = mocks.NewMockPostService(mockCtrl)
	mockRenderService = mocks.NewMockRenderService(mockCtrl)
	mockPostCollaboratorService = mocks.NewMockPostCollaboratorService(mockCtrl)

	// Setup SUT
	postController = PostController{
		PostService:             mockPostService,
		RenderService:           mockRenderService,
		PostCollaboratorService: mockPostCollaboratorService,
	}

	// Setup HTTP response recorder
	responseRecorder = httptest.NewRecorder()
}

func teardownPostController() {
	_ = lock.Unlock()
}

// Helper function that creates a multi-part form data body to send in a HTTP request
// Returns the body, and the form data content type
func createTestingFormDataBody() (*bytes.Buffer, string, error) {
	// Setup buffer to write file to
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Attach writer to file
	part, err := writer.CreateFormFile("file", "test.zip")
	if err != nil {
		return nil, "", err
	}

	// Write content to file
	if _, err := part.Write([]byte("file content")); err != nil {
		return nil, "", err
	}

	writer.Close()

	return body, writer.FormDataContentType(), nil
}

func TestGetAllPostCollaboratorsGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	postID := uint(10)

	post := &models.Post{
		Model: gorm.Model{ID: postID},
		Collaborators: []*models.PostCollaborator{
			{
				Model:             gorm.Model{ID: 5},
				Member:            models.Member{Model: gorm.Model{ID: 2}},
				MemberID:          2,
				PostID:            postID,
				CollaborationType: models.Author,
			},

			{
				Model:             gorm.Model{ID: 12},
				Member:            models.Member{Model: gorm.Model{ID: 3}},
				MemberID:          3,
				PostID:            postID,
				CollaborationType: models.Contributor,
			},
			{
				Model:             gorm.Model{ID: 15},
				Member:            models.Member{Model: gorm.Model{ID: 10}},
				MemberID:          10,
				PostID:            postID,
				CollaborationType: models.Reviewer,
			},
		},
	}

	// Setup mocks
	mockPostService.EXPECT().GetPost(postID).Return(post, nil).Times(1)

	// Construct req
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/collaborators/all/%d", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Read response
	responseJSON, err := io.ReadAll(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Parse response
	responsePostCollaboratorDTOs := []*models.PostCollaboratorDTO{}
	if err := json.Unmarshal(responseJSON, &responsePostCollaboratorDTOs); err != nil {
		t.Fatal(err)
	}

	expectedPostCollaboratorDTOs := []*models.PostCollaboratorDTO{
		{
			ID:                5,
			MemberID:          2,
			PostID:            postID,
			CollaborationType: models.Author,
		},
		{
			ID:                12,
			MemberID:          3,
			PostID:            postID,
			CollaborationType: models.Contributor,
		},
		{
			ID:                15,
			MemberID:          10,
			PostID:            postID,
			CollaborationType: models.Reviewer,
		},
	}

	assert.Equal(t, expectedPostCollaboratorDTOs, responsePostCollaboratorDTOs)
}

func TestGetAllPostCollaboratorsPostNotFound(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	postID := uint(10)

	// Setup mocks
	mockPostService.EXPECT().GetPost(postID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/collaborators/all/%d", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetAllPostCollaboratorsInvalidPostID(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Construct request
	req, err := http.NewRequest("GET", "/api/v2/posts/collaborators/all/badID", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetPostGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(99)
	collaboratorID := uint(5)
	memberID := uint(1)
	tagContainerID := uint(5)
	discussionContainerID := uint(10)
	createdAt := time.Now().UTC()
	updatedAt := time.Now().Add(time.Minute).UTC()

	post := &models.Post{
		Model: gorm.Model{ID: postID, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Collaborators: []*models.PostCollaborator{
			{
				Model:             gorm.Model{ID: collaboratorID},
				Member:            models.Member{},
				MemberID:          memberID,
				PostID:            postID,
				CollaborationType: models.Author,
			},
		},
		Title:    "my cool post",
		PostType: models.Question,
		ScientificFieldTagContainer: models.ScientificFieldTagContainer{
			Model:               gorm.Model{ID: tagContainerID},
			ScientificFieldTags: []*models.ScientificFieldTag{},
		},
		ScientificFieldTagContainerID: tagContainerID,
		DiscussionContainer: models.DiscussionContainer{
			Model:       gorm.Model{ID: discussionContainerID},
			Discussions: []*models.Discussion{},
		},
		DiscussionContainerID: discussionContainerID,
		RenderStatus:          models.Pending,
	}

	// Setup mocks
	mockPostService.EXPECT().GetPost(postID).Return(post, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseBody := &models.PostDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseBody); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedBody := &models.PostDTO{
		ID:                            postID,
		CollaboratorIDs:               []uint{collaboratorID},
		Title:                         "my cool post",
		PostType:                      models.Question,
		ScientificFieldTagContainerID: tagContainerID,
		DiscussionContainerID:         discussionContainerID,
		RenderStatus:                  models.Pending,
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}

	assert.Equal(t, expectedBody, responseBody)
}

func TestGetPostNotFound(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(5)

	// Setup mocks
	mockPostService.EXPECT().GetPost(postID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreatePostGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	form := &forms.PostCreationForm{
		AuthorMemberIDs:       []uint{},
		Title:                 "my post",
		Anonymous:             true,
		PostType:              models.Question,
		ScientificFieldTagIDs: []uint{},
	}

	postID := uint(5)
	tagContainerID := uint(10)
	discussionContainerID := uint(15)
	createdAt := time.Now().UTC()
	updatedAt := time.Now().Add(time.Minute).UTC()

	// Setup mocks
	mockPostService.EXPECT().CreatePost(form).Return(&models.Post{
		Model:         gorm.Model{ID: postID, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Collaborators: []*models.PostCollaborator{},
		Title:         "my post",
		PostType:      models.Question,
		ScientificFieldTagContainer: models.ScientificFieldTagContainer{
			Model:               gorm.Model{ID: tagContainerID},
			ScientificFieldTags: []*models.ScientificFieldTag{},
		},
		ScientificFieldTagContainerID: tagContainerID,
		DiscussionContainer: models.DiscussionContainer{
			Model:       gorm.Model{ID: discussionContainerID},
			Discussions: []*models.Discussion{},
		},
		DiscussionContainerID: discussionContainerID,
		RenderStatus:          models.Pending,
	}, nil).Times(1)

	// Marshal form
	body, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/posts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseBody := &models.PostDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseBody); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedBody := &models.PostDTO{
		ID:                            postID,
		CollaboratorIDs:               []uint{},
		Title:                         "my post",
		PostType:                      models.Question,
		ScientificFieldTagContainerID: tagContainerID,
		DiscussionContainerID:         discussionContainerID,
		RenderStatus:                  models.Pending,
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}

	assert.Equal(t, expectedBody, responseBody)
}

func TestCreatePostFormValidationFailed(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	form := forms.PostCreationForm{
		AuthorMemberIDs:       []uint{1, 2},
		Title:                 "",
		Anonymous:             false,
		PostType:              "",
		ScientificFieldTagIDs: []uint{1},
	}

	// Marshal form
	body, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/posts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreatePostDatabaseCreationFailed(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	form := forms.PostCreationForm{
		AuthorMemberIDs:       []uint{1, 2},
		Title:                 "my awesome post",
		Anonymous:             false,
		PostType:              models.Question,
		ScientificFieldTagIDs: []uint{1},
	}

	// Setup mocks
	mockPostService.EXPECT().CreatePost(&form).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Marshal form
	body, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/posts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetPostCollaboratorGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	collaboratorID := uint(5)
	memberID := uint(10)
	postID := uint(1)

	collaborator := &models.PostCollaborator{
		Model:             gorm.Model{ID: collaboratorID},
		Member:            models.Member{Model: gorm.Model{ID: memberID}},
		MemberID:          memberID,
		PostID:            postID,
		CollaborationType: models.Author,
	}

	// Setup mocks
	mockPostCollaboratorService.EXPECT().GetPostCollaborator(collaboratorID).Return(collaborator, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/collaborators/%d", collaboratorID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseDTO := &models.PostCollaboratorDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedDTO := &models.PostCollaboratorDTO{
		ID:                collaboratorID,
		MemberID:          memberID,
		PostID:            postID,
		CollaborationType: models.Author,
	}

	assert.Equal(t, expectedDTO, responseDTO)
}

func TestGetPostCollaboratorNotFound(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	collaboratorID := uint(5)

	// Setup mocks
	mockPostCollaboratorService.EXPECT().GetPostCollaborator(collaboratorID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/collaborators/%d", collaboratorID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetPostCollaboratorBadID(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	collaboratorID := "bad id"

	// Setup mocks
	mockPostCollaboratorService.EXPECT().GetPostCollaborator(collaboratorID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/collaborators/%s", collaboratorID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUploadPostGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(5)

	body, contentType, err := createTestingFormDataBody()
	if err != nil {
		t.Fatal(err)
	}

	// Setup mocks
	mockPostService.EXPECT().UploadPost(gomock.Any(), gomock.Any(), postID).Return(nil).Times(1)

	// Construct request
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/posts/%d/upload", postID), body)
	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUploadPostNoFilesAttached(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Construct request
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/posts/%d/upload", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUploadPostUploadingFailed(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	body, contentType, err := createTestingFormDataBody()
	if err != nil {
		t.Fatal(err)
	}

	// Setup mocks
	mockPostService.EXPECT().UploadPost(gomock.Any(), gomock.Any(), postID).Return(fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/posts/%d/upload", postID), body)
	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestUploadPostInvalidPostID(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := "bad id!!!!!"

	body, contentType, err := createTestingFormDataBody()
	if err != nil {
		t.Fatal(err)
	}

	// Setup mocks
	mockPostService.EXPECT().UploadPost(gomock.Any(), gomock.Any(), postID).Return(fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/posts/%s/upload", postID), body)
	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetMainRenderGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)
	filePath := "../utils/test_files/good_repository_setup/render/1234.html"

	// Setup mocks
	mockRenderService.EXPECT().GetMainRenderFile(postID).Return(filePath, lock, nil, nil, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/render", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "text/html", responseRecorder.Header().Get("Content-Type"))
	assert.False(t, lock.Locked())
}

func TestGetMainRenderInvalidPostID(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := "bad id!!"

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%s/render", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetMainRenderPending(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Setup mocks
	mockRenderService.EXPECT().GetMainRenderFile(postID).Return("", nil, fmt.Errorf("oh no"), nil, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/render", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusAccepted, responseRecorder.Result().StatusCode)
}

func TestGetMainRenderFailed(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Setup mocks
	mockRenderService.EXPECT().GetMainRenderFile(postID).Return("", nil, nil, fmt.Errorf("oh no"), nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/render", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNoContent, responseRecorder.Result().StatusCode)
}

func TestGetMainFiletreeGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Setup mocks
	mockPostService.EXPECT().GetMainFiletree(postID).Return(map[string]int64{
		"./my_file":   42,
		"./my_folder": -1,
	}, nil, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/tree", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	response := &map[string]int64{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(response); err != nil {
		t.Fatal(err)
	}

	// Check body
	expected := &map[string]int64{
		"./my_file":   42,
		"./my_folder": -1,
	}

	assert.Equal(t, expected, response)
}

func TestGetMainFileTreePostNotFound(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Setup mocks
	mockPostService.EXPECT().GetMainFiletree(postID).Return(nil, fmt.Errorf("oh no"), nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/tree", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetMainFileTreeGetTreeFailed(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Setup mocks
	mockPostService.EXPECT().GetMainFiletree(postID).Return(nil, nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/tree", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostIfExistsGoodWeather(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)
	projectPostID := uint(15)

	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: projectPostID},
		PostID: postID,
	}

	// Setup mocks
	mockPostService.EXPECT().GetProjectPost(postID).Return(projectPost, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/project-post", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	response := new(uint)
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(response); err != nil {
		t.Fatal(err)
	}

	// Check body
	expected := projectPostID

	assert.Equal(t, expected, *response)
}

func TestGetProjectPostIfExistsPostDNE(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)

	// Setup mocks
	mockPostService.EXPECT().GetProjectPost(postID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/project-post", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetProjectPostIfExistsBadPostID(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := "bad id!!"

	// Setup mocks
	mockPostService.EXPECT().GetProjectPost(postID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%s/project-post", postID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetMainFileFromProject(t *testing.T) {
	setupPostController(t)
	t.Cleanup(teardownPostController)

	// Setup data
	postID := uint(10)
	relativeFilePath := "./1234.html"
	absoluteFilePath := "../utils/test_files/good_repository_setup/render/1234.html"

	// Setup mocks
	mockPostService.EXPECT().GetMainFileFromProject(postID, gomock.Any()).Return(absoluteFilePath, lock, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/posts/%d/file/%s", postID, relativeFilePath), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.False(t, lock.Locked())
}
