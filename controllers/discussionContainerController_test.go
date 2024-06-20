package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func setupDiscussionContainerController(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockDiscussionContainerService = mocks.NewMockDiscussionContainerService(mockCtrl)

	// Setup SUT
	discussionContainerController = DiscussionContainerController{
		DiscussionContainerService: mockDiscussionContainerService,
	}

	// Setup HTTP
	responseRecorder = httptest.NewRecorder()
}

func teardownDiscussionContainerController() {

}

func TestGetDiscussionContainerGoodWeather(t *testing.T) {
	setupDiscussionContainerController(t)
	t.Cleanup(teardownDiscussionContainerController)

	// Setup data
	discussionContainerID := uint(10)

	discussionContainer := &models.DiscussionContainer{
		Model: gorm.Model{ID: discussionContainerID},
		Discussions: []*models.Discussion{
			{Model: gorm.Model{ID: 20}},
			{Model: gorm.Model{ID: 25}},
		},
	}

	// Setup mocks
	mockDiscussionContainerService.EXPECT().GetDiscussionContainer(discussionContainerID).Return(discussionContainer, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/discussion-containers/%d", discussionContainerID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseDiscussionContainerDTO := &models.DiscussionContainerDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseDiscussionContainerDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedDiscussionContainerDTO := &models.DiscussionContainerDTO{
		ID:            discussionContainerID,
		DiscussionIDs: []uint{20, 25},
	}

	assert.Equal(t, expectedDiscussionContainerDTO, responseDiscussionContainerDTO)
}

func TestGetDiscussionContainerMalformedID(t *testing.T) {
	setupDiscussionContainerController(t)
	t.Cleanup(teardownDiscussionContainerController)

	// Setup data
	discussionContainerID := "bad!!!"

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/discussion-containers/%s", discussionContainerID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetDiscussionContainerNotFound(t *testing.T) {
	setupDiscussionContainerController(t)
	t.Cleanup(teardownDiscussionContainerController)

	// Setup data
	discussionContainerID := uint(10)

	// Setup mocks
	mockDiscussionContainerService.EXPECT().GetDiscussionContainer(discussionContainerID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/discussion-containers/%d", discussionContainerID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}
