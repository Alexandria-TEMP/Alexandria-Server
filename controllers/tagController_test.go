package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func beforeEachTag(t *testing.T) {
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	// Setup mocks
	mockTagService = mock_interfaces.NewMockTagService(mockCtrl)
	mockScientificFieldTagContainerService = mock_interfaces.NewMockScientificFieldTagContainerService(mockCtrl)

	// Setup SUT
	tagController = TagController{
		TagService:                         mockTagService,
		ScientificFieldTagContainerService: mockScientificFieldTagContainerService,
	}

	// Setup HTTP testing
	responseRecorder = httptest.NewRecorder()
}

func TestGetScientificTags200(t *testing.T) {
	// call the before each method
	beforeEachTag(t)

	// define response from mock tag service
	mockTagService.EXPECT().GetAllScientificFieldTags().Return([]*models.ScientificFieldTag{exampleSTag1, exampleSTag2}, nil)

	// set up request
	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	// set up receiving a response
	var responseTags []models.ScientificFieldTagDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responseTags)

	// set up expected response
	expectedResponseTags := []models.ScientificFieldTagDTO{exampleSTag1.IntoDTO(), exampleSTag2.IntoDTO()}

	// assert that response is as expected
	assert.Equal(t, expectedResponseTags, responseTags)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetScientificTags404(t *testing.T) {
	beforeEachTag(t)

	// set up mock tag service to return an error
	mockTagService.EXPECT().GetAllScientificFieldTags().Return(nil, errors.New("some error")).Times(1)

	// set up request
	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	// assert that the correct response was returned
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetScientificFieldTag200(t *testing.T) {
	beforeEachTag(t)

	mockTagService.EXPECT().GetTagByID(uint(1)).Return(exampleSTag1, nil).Times(1)

	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsetag models.ScientificFieldTagDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsetag)

	assert.Equal(t, exampleSTag1DTO, responsetag)
}

func TestGetScientificFieldTag400(t *testing.T) {
	beforeEachTag(t)

	mockTagService.EXPECT().GetTagByID(gomock.Any()).Return(exampleSTag1, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetScientificFieldTag404(t *testing.T) {
	beforeEachTag(t)

	mockTagService.EXPECT().GetTagByID(uint(1)).Return(&models.ScientificFieldTag{}, errors.New("some error")).Times(1)

	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetScientificFieldTagContainerGoodWeather(t *testing.T) {
	beforeEachTag(t)

	containerID := uint(20)

	container := &models.ScientificFieldTagContainer{
		Model: gorm.Model{ID: containerID},
		ScientificFieldTags: []*models.ScientificFieldTag{
			{Model: gorm.Model{ID: 2}},
			{Model: gorm.Model{ID: 3}},
			{Model: gorm.Model{ID: 5}},
		},
	}

	// Setup mocks
	mockScientificFieldTagContainerService.EXPECT().GetScientificFieldTagContainer(containerID).Return(container, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/tags/scientific/containers/%d", containerID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check response
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Read body
	bytes, err := io.ReadAll(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Parse body
	var responseScientificFieldTagContainerDTO models.ScientificFieldTagContainerDTO
	if err := json.Unmarshal(bytes, &responseScientificFieldTagContainerDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedScientificFieldTagContainerDTO := models.ScientificFieldTagContainerDTO{
		ID:                    containerID,
		ScientificFieldTagIDs: []uint{2, 3, 5},
	}

	assert.Equal(t, expectedScientificFieldTagContainerDTO, responseScientificFieldTagContainerDTO)
}

func TestGetCompletionStatusTags(t *testing.T) {
	beforeEachTag(t)

	// Construct request
	req, err := http.NewRequest("GET", "/api/v2/tags/completion-status", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check response
	// We do NOT test the BODY return value of this endpoint - it's 100% hard-coded,
	// and therefore testing it would only serve to make it brittle & resistant to change.
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetPostTypeTags(t *testing.T) {
	beforeEachTag(t)

	// Construct request
	req, err := http.NewRequest("GET", "/api/v2/tags/post-type", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check response
	// We do NOT test the BODY return value of this endpoint - it's 100% hard-coded,
	// and therefore testing it would only serve to make it brittle & resistant to change.
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestGetFeedbackPreferenceTags(t *testing.T) {
	beforeEachTag(t)

	// Construct request
	req, err := http.NewRequest("GET", "/api/v2/tags/feedback-preference", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check response
	// We do NOT test the BODY return value of this endpoint - it's 100% hard-coded,
	// and therefore testing it would only serve to make it brittle & resistant to change.
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}
