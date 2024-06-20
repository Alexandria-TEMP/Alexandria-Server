package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func setupDiscussionController(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockDiscussionService = mocks.NewMockDiscussionService(mockCtrl)

	// Setup SUT
	discussionController = DiscussionController{
		DiscussionService: mockDiscussionService,
	}

	// Setup HTTP
	responseRecorder = httptest.NewRecorder()
}

func teardownDiscussionController() {

}

func TestGetDiscussionGoodWeather(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	discussionID := uint(10)
	memberID := uint(5)
	parentID := uint(7)
	replyID := uint(20)
	containerID := uint(1)

	discussion := &models.Discussion{
		Model:       gorm.Model{ID: discussionID},
		ContainerID: containerID,
		Member:      &models.Member{Model: gorm.Model{ID: memberID}},
		MemberID:    &memberID,
		Replies: []*models.Discussion{
			{
				Model:       gorm.Model{ID: replyID},
				ContainerID: containerID,
				Member:      &models.Member{Model: gorm.Model{ID: memberID}},
				MemberID:    &memberID,
				Replies:     []*models.Discussion{},
				ParentID:    &discussionID,
				Text:        "my reply",
			},
		},
		ParentID: &parentID,
		Text:     "my discussion",
	}

	// Setup mocks
	mockDiscussionService.EXPECT().GetDiscussion(discussionID).Return(discussion, nil).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/discussions/%d", discussionID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseDiscussionDTO := &models.DiscussionDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseDiscussionDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedDiscussionDTO := &models.DiscussionDTO{
		ID:       discussionID,
		MemberID: &memberID,
		ReplyIDs: []uint{replyID},
		Text:     "my discussion",
	}

	assert.Equal(t, expectedDiscussionDTO, responseDiscussionDTO)
}

func TestGetDiscussionInvalidID(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	discussionID := "bad id!!!"

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/discussions/%s", discussionID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetDiscussionNotFound(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	discussionID := uint(10)

	// Setup mocks
	mockDiscussionService.EXPECT().GetDiscussion(discussionID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/discussions/%d", discussionID), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestCreateRootDiscussionGoodWeather(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	containerID := uint(10)
	memberID := uint(5)

	rootDiscussionCreationForm := forms.RootDiscussionCreationForm{
		ContainerID: containerID,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "my root discussion",
		},
	}
	body, _ := json.Marshal(rootDiscussionCreationForm)
	member := &models.Member{Model: gorm.Model{ID: memberID}}

	// Setup mocks
	mockDiscussionService.EXPECT().CreateRootDiscussion(&rootDiscussionCreationForm, member).Return(&models.Discussion{
		Model:       gorm.Model{ID: 5},
		ContainerID: containerID,
		Member:      member,
		MemberID:    &memberID,
		Replies:     []*models.Discussion{},
		ParentID:    nil,
		Text:        "my root discussion",
	}, nil).Times(1)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	discussionController.CreateRootDiscussion(c)

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseDiscussionDTO := &models.DiscussionDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseDiscussionDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedDiscussionDTO := &models.DiscussionDTO{
		ID:       5,
		MemberID: &memberID,
		ReplyIDs: []uint{},
		Text:     "my root discussion",
	}

	assert.Equal(t, expectedDiscussionDTO, responseDiscussionDTO)
}

func TestCreateRootDiscussionInvalidForm(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	rootDiscussionCreationForm := &forms.RootDiscussionCreationForm{
		ContainerID: 0,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "", // An empty discussion will fail validation
		},
	}

	// Marshal form
	body, err := json.Marshal(rootDiscussionCreationForm)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("POST", "/api/v2/discussions/roots", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateRootDiscussionDatabaseFailure(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	rootDiscussionCreationForm := &forms.RootDiscussionCreationForm{
		ContainerID: 0,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "my discussion",
		},
	}
	body, _ := json.Marshal(rootDiscussionCreationForm)
	member := &models.Member{}

	// Setup mocks
	mockDiscussionService.EXPECT().CreateRootDiscussion(rootDiscussionCreationForm, member).Return(nil, fmt.Errorf("oh no")).Times(1)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	discussionController.CreateRootDiscussion(c)

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestCreateReplyDiscussionGoodWeather(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	parentID := uint(10)
	containerID := uint(7)
	memberID := uint(5)

	replyDiscussionCreationForm := forms.ReplyDiscussionCreationForm{
		ParentID: parentID,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "my reply discussion",
		},
	}
	body, _ := json.Marshal(replyDiscussionCreationForm)
	member := &models.Member{Model: gorm.Model{ID: memberID}}

	// Setup mocks
	mockDiscussionService.EXPECT().CreateReply(&replyDiscussionCreationForm, member).Return(&models.Discussion{
		Model:       gorm.Model{ID: 5},
		ContainerID: containerID,
		Member:      &models.Member{},
		MemberID:    &memberID,
		Replies:     []*models.Discussion{},
		ParentID:    &parentID,
		Text:        "my reply discussion",
	}, nil).Times(1)

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	discussionController.CreateReplyDiscussion(c)

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseDiscussionDTO := &models.DiscussionDTO{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseDiscussionDTO); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedDiscussionDTO := &models.DiscussionDTO{
		ID:       5,
		MemberID: &memberID,
		ReplyIDs: []uint{},
		Text:     "my reply discussion",
	}

	assert.Equal(t, expectedDiscussionDTO, responseDiscussionDTO)
}

func TestCreateReplyDiscussionInvalidForm(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	replyDiscussionCreationForm := &forms.ReplyDiscussionCreationForm{
		ParentID: 0,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "", // An empty discussion will fail validation
		},
	}
	body, _ := json.Marshal(replyDiscussionCreationForm)
	member := &models.Member{}

	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	discussionController.CreateReplyDiscussion(c)

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateReplyDiscussionDatabaseFailure(t *testing.T) {
	setupDiscussionController(t)
	t.Cleanup(teardownDiscussionController)

	// Setup data
	replyDiscussionCreationForm := &forms.ReplyDiscussionCreationForm{
		ParentID: 0,
		DiscussionCreationForm: forms.DiscussionCreationForm{
			Anonymous: false,
			Text:      "my discussion",
		},
	}
	body, _ := json.Marshal(replyDiscussionCreationForm)
	member := &models.Member{}

	// Setup mocks
	mockDiscussionService.EXPECT().CreateReply(replyDiscussionCreationForm, member).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Marshal form
	c, _ := gin.CreateTestContext(responseRecorder)
	c.Set("currentMember", member)
	c.Request = &http.Request{}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	discussionController.CreateReplyDiscussion(c)

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}
