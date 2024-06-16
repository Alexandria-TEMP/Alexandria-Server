package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
