package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

<<<<<<<< HEAD:tests/controllerTests/postController_test.go
========
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
>>>>>>>> main:controllers/postController_test.go
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

<<<<<<<< HEAD:tests/controllerTests/postController_test.go
func beforeEachController(t *testing.T) {
========
var (
	mockPostService *mock_interfaces.MockPostService
	postController  *PostController
	router          *gin.Engine

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
	router.PUT("/api/v1/post", func(c *gin.Context) {
		postController.UpdatePost(c)
	})
	router.GET("/api/v1/projectPost/:postID", func(c *gin.Context) {
		postController.GetProjectPost(c)
	})
	router.POST("/api/v1/projectPost", func(c *gin.Context) {
		postController.CreateProjectPost(c)
	})
	router.PUT("/api/v1/projectPost", func(c *gin.Context) {
		postController.UpdateProjectPost(c)
	})

	// Setup object
	examplePost = models.Post{Model: gorm.Model{ID: 1}}
	exampleProjectPost = models.ProjectPost{Model: gorm.Model{ID: 2}}

	os.Exit(m.Run())
}

func beforeEach(t *testing.T) {
>>>>>>>> main:controllers/postController_test.go
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	responseRecorder = httptest.NewRecorder()

	mockPostService = mock_interfaces.NewMockPostService(mockCtrl)
	postController = &PostController{PostService: mockPostService}
}

func TestGetPost200(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().GetPost(uint64(1)).Return(&examplePost, nil).Times(1)

	req, _ := http.NewRequest("GET", "/api/v1/post/1", http.NoBody)

	router.ServeHTTP(responseRecorder, req)

	var responsePost models.Post

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsePost)

	assert.Equal(t, examplePost, responsePost)
}

func TestGetPost400(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().GetPost(gomock.Any()).Return(&examplePost, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/post/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetPost410(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().GetPost(uint64(1)).Return(&models.Post{}, errors.New("some error")).Times(1)

	examplePostJSON, _ := json.Marshal(examplePost)
	req, _ := http.NewRequest("GET", "/api/v1/post/1", bytes.NewBuffer(examplePostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}

func TestCreatePost200(t *testing.T) {
	beforeEachController(t)

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
	beforeEachController(t)

	mockPostService.EXPECT().CreatePost(gomock.Any()).Return(&examplePost).Times(0)

	badPostFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(badPostFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdatePost200(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().UpdatePost(&examplePost).Return(nil).Times(1)

	examplePostJSON, _ := json.Marshal(examplePost)
	req, _ := http.NewRequest("PUT", "/api/v1/post", bytes.NewBuffer(examplePostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUpdatePost400(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().UpdatePost(&examplePost).Return(nil).Times(0)

	examplePostJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("PUT", "/api/v1/post", bytes.NewBuffer(examplePostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdatePost410(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().UpdatePost(&examplePost).Return(errors.New("some error")).Times(1)

	examplePostJSON, _ := json.Marshal(examplePost)
	req, _ := http.NewRequest("PUT", "/api/v1/post", bytes.NewBuffer(examplePostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}

func TestGetProjectPost200(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().GetProjectPost(uint64(2)).Return(&exampleProjectPost, nil).Times(1)

	req, _ := http.NewRequest("GET", "/api/v1/projectPost/2", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	var responsePost models.ProjectPost

	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &responsePost)

	assert.Equal(t, exampleProjectPost, responsePost)
}

// func TestGetProjectPostDoesntExist(t *testing.T) {
// 	beforeEachController(t)
// }

func TestGetProjectPost400(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().GetProjectPost(gomock.Any()).Return(&exampleProjectPost, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/projectPost/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetProjectPost410(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().GetProjectPost(uint64(1)).Return(&models.ProjectPost{}, errors.New("some error")).Times(1)

	exampleProjectPostJSON, _ := json.Marshal(exampleProjectPost)
	req, _ := http.NewRequest("GET", "/api/v1/projectPost/1", bytes.NewBuffer(exampleProjectPostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}

func TestCreateProjectPost200(t *testing.T) {
	beforeEachController(t)

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
	beforeEachController(t)

	mockPostService.EXPECT().CreateProjectPost(gomock.Any()).Return(&exampleProjectPost).Times(0)

	badProjectPostFormJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("POST", "/api/v1/projectPost", bytes.NewBuffer(badProjectPostFormJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateProjectPost200(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().UpdateProjectPost(&exampleProjectPost).Return(nil).Times(1)

	exampleProjectPostJSON, _ := json.Marshal(exampleProjectPost)
	req, _ := http.NewRequest("PUT", "/api/v1/projectPost", bytes.NewBuffer(exampleProjectPostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestUpdateProjectPost400(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().UpdateProjectPost(gomock.Any()).Return(nil).Times(0)

	exampleProjectPostJSON := []byte(`jgdfskljglkdjlmdflkgmlksdfglksdlfgdsfgsdg`)
	req, _ := http.NewRequest("PUT", "/api/v1/projectPost", bytes.NewBuffer(exampleProjectPostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestUpdateProjectPost410(t *testing.T) {
	beforeEachController(t)

	mockPostService.EXPECT().UpdateProjectPost(&exampleProjectPost).Return(errors.New("some error")).Times(1)

	exampleProjectPostJSON, _ := json.Marshal(exampleProjectPost)
	req, _ := http.NewRequest("PUT", "/api/v1/projectPost", bytes.NewBuffer(exampleProjectPostJSON))
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusGone, responseRecorder.Result().StatusCode)
}
