package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"go.uber.org/mock/gomock"
)

func setupFilterController(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockPostService = mocks.NewMockPostService(mockCtrl)

	// Setup SUT
	filterController = FilterController{
		PostService: mockPostService,
	}

	// Setup HTTP
	responseRecorder = httptest.NewRecorder()
}

func teardownFilterController() {

}

func TestFilterPostsGoodWeather(t *testing.T) {
	setupFilterController(t)
	t.Cleanup(teardownFilterController)

	// Setup data
	filterForm := forms.PostFilterForm{
		IncludeProjectPosts: true,
	}

	page := 1
	size := 10

	postID1 := uint(10)
	postID2 := uint(15)

	// Setup mocks
	mockPostService.EXPECT().Filter(page, size, filterForm).Return([]uint{postID1, postID2}, nil).Times(1)

	// Marshal form
	body, err := json.Marshal(filterForm)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/filter/posts?page=%d&size=%d", page, size), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	// Decode body
	responseBody := &[]uint{}
	if err := json.NewDecoder(responseRecorder.Result().Body).Decode(responseBody); err != nil {
		t.Fatal(err)
	}

	// Check body
	expectedBody := &[]uint{postID1, postID2}

	assert.Equal(t, expectedBody, responseBody)
}

func TestFilterPostsDatabaseQueryFailed(t *testing.T) {
	setupFilterController(t)
	t.Cleanup(teardownFilterController)

	// Setup data
	filterForm := forms.PostFilterForm{
		IncludeProjectPosts: false,
	}

	page := 1
	size := 10

	// Setup mocks
	mockPostService.EXPECT().Filter(page, size, filterForm).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Marshal form
	body, err := json.Marshal(filterForm)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/filter/posts?page=%d&size=%d", page, size), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestFilterPostsFormBadlyFormatted(t *testing.T) {
	setupFilterController(t)
	t.Cleanup(teardownFilterController)

	// Setup data
	type BadPostFilterForm struct {
		IncludeProjectPosts string
	}

	filterForm := BadPostFilterForm{
		IncludeProjectPosts: "ooga",
	}

	page := 1
	size := 10

	// Marshal form
	body, err := json.Marshal(filterForm)
	if err != nil {
		t.Fatal(err)
	}

	// Construct request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v2/filter/posts?page=%d&size=%d", page, size), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Send request
	router.ServeHTTP(responseRecorder, req)
	defer responseRecorder.Result().Body.Close()

	// Check status
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}
