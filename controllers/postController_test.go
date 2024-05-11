package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

var (
	mockPostService        *mock_interfaces.MockPostService
	postController         *PostController
	router                 *gin.Engine
	responseRecorder       *httptest.ResponseRecorder
	examplePost            models.Post
	exampleProjectPost     models.ProjectPost
	examplePostForm        forms.PostCreationForm
	exampleProjectPostForm forms.ProjectPostCreationForm
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = gin.Default()
	gin.SetMode(gin.TestMode)

	router.GET("/api/v1/post/:postID", func(c *gin.Context) {
		postController.GetPost(c)
	})
	router.POST("/api/v1/post", func(c *gin.Context) {
		postController.CreatePost(c)
	})
	router.GET("/api/v1/projectPost/:postID", func(c *gin.Context) {
		postController.GetProjectPost(c)
	})
	router.POST("/api/v1/projectPost", func(c *gin.Context) {
		postController.CreateProjectPost(c)
	})

	// Setup object
	examplePost = models.Post{ID: 1}
	exampleProjectPost = models.ProjectPost{ID: 2}

	os.Exit(m.Run())
}

func beforeEach(t *testing.T) {
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	responseRecorder = httptest.NewRecorder()

	mockPostService = mock_interfaces.NewMockPostService(mockCtrl)
	postController = &PostController{PostService: mockPostService}
}

func TestGetPost200(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().GetPost(uint64(1)).Return(&examplePost).Times(1)

	req, _ := http.NewRequest("GET", "/api/v1/post/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsePost models.Post

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsePost)

	assert.Equal(t, examplePost, responsePost)
}

func TestGetPostDoesntExist(t *testing.T) {
	beforeEach(t)
}

func TestGetPost400(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().GetPost(gomock.Any()).Return(&examplePost).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/post/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreatePost200(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().CreatePost(&examplePostForm).Return(&examplePost).Times(1)

	examplePostFormJSON, _ := json.Marshal(examplePostForm)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(examplePostFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsePost models.Post

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsePost)

	assert.Equal(t, examplePost, responsePost)
}

func TestCreatePost400(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().CreatePost(gomock.Any()).Return(&examplePost).Times(0)

	badPostFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(badPostFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetProjectPost200(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().GetProjectPost(uint64(2)).Return(&exampleProjectPost).Times(1)

	req, _ := http.NewRequest("GET", "/api/v1/projectPost/2", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsePost models.ProjectPost

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsePost)

	assert.Equal(t, exampleProjectPost, responsePost)
}

// func TestGetProjectPostDoesntExist(t *testing.T) {
// 	beforeEach(t)
// }

func TestGetProjectPost400(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().GetProjectPost(gomock.Any()).Return(&exampleProjectPost).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/projectPost/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateProjectPost200(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().CreateProjectPost(&exampleProjectPostForm).Return(&exampleProjectPost).Times(1)

	exampleProjectPostFormJSON, _ := json.Marshal(exampleProjectPostForm)
	req, _ := http.NewRequest("POST", "/api/v1/projectPost", bytes.NewBuffer(exampleProjectPostFormJSON))
	router.ServeHTTP(responseRecorder, req)

	var responsePost models.ProjectPost

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsePost)

	assert.Equal(t, exampleProjectPost, responsePost)
}

func TestCreateProjectPost400(t *testing.T) {
	beforeEach(t)

	mockPostService.EXPECT().CreateProjectPost(gomock.Any()).Return(&exampleProjectPost).Times(0)

	badProjectPostFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v1/projectPost", bytes.NewBuffer(badProjectPostFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}
