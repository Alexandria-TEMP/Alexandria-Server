package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	gomock "go.uber.org/mock/gomock"
)

func beforeEachTag(t *testing.T) {
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	responseRecorder = httptest.NewRecorder()

	mockTagService = mock_interfaces.NewMockTagService(mockCtrl)
	tagController = &TagController{TagService: mockTagService}
}

func TestGetScientificTags200(t *testing.T) {
	// call the before each method
	beforeEachTag(t)

	// define response from mock tag service
	mockTagService.EXPECT().GetAllScientificFieldTags().Return([]*tags.ScientificFieldTag{exampleSTag1, exampleSTag2}, nil)

	// set up request
	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	// set up receiving a response
	var responseTags []tags.ScientificFieldTagDTO

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responseTags)

	// set up expected response
	expectedResponseTags := []tags.ScientificFieldTagDTO{exampleSTag1.IntoDTO(), exampleSTag2.IntoDTO()}

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

	var responsetag tags.ScientificFieldTagDTO

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

	mockTagService.EXPECT().GetTagByID(uint(1)).Return(&tags.ScientificFieldTag{}, errors.New("some error")).Times(1)

	req, _ := http.NewRequest("GET", "/api/v2/tags/scientific/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}
