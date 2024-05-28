package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"go.uber.org/mock/gomock"
)

func beforeEachVersion(t *testing.T) {
	t.Helper()
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockVersionService = mock_interfaces.NewMockVersionService(mockCtrl)
	versionController = &VersionController{VersionService: mockVersionService}

	responseRecorder = httptest.NewRecorder()
}

func TestCreateVersion200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any()).Return(&examplePendingVersion, nil).Times(1)

	zipPath := filepath.Join(cwd, "..", "utils", "test_files", "file_handling_test.zip")
	body, dataType, err := CreateMultipartFile(zipPath)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/versions", body)
	req.Header.Add("Content-Type", dataType)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCreateVersion4001(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any()).Return(&examplePendingVersion, nil).Times(0)

	req, _ := http.NewRequest("POST", "/api/v2/versions", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateVersion500(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any()).Return(&examplePendingVersion, errors.New("err")).Times(1)

	zipPath := filepath.Join(cwd, "..", "utils", "test_files", "file_handling_test.zip")
	body, dataType, err := CreateMultipartFile(zipPath)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v2/versions", body)
	req.Header.Add("Content-Type", dataType)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetRender200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(uint(1)).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil)

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "attachment; filename=render.html", responseRecorder.Result().Header.Get("Content-Disposition"))

	fileHeader := make([]byte, headerSize)
	_, err := responseRecorder.Result().Body.Read(fileHeader)
	assert.Nil(t, err)

	fileContentType := http.DetectContentType(fileHeader)

	assert.Equal(t, fileContentType, "text/html; charset=utf-8")
}

func TestGetRender202(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(uint(1)).Return("", errors.New("err"), nil)

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusAccepted, responseRecorder.Result().StatusCode)
}

func TestGetRender4001(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(gomock.Any()).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/bad/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender4002(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(gomock.Any()).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/-1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender404(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(uint(1)).Return("", nil, errors.New("err"))

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetReposiotry200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(uint(1)).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil)

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "attachment; filename=quarto_project.zip", responseRecorder.Result().Header.Get("Content-Disposition"))

	fileHeader := make([]byte, headerSize)
	_, err := responseRecorder.Result().Body.Read(fileHeader)
	assert.Nil(t, err)

	fileContentType := http.DetectContentType(fileHeader)

	assert.Equal(t, fileContentType, "application/zip")
}

func TestGetReposiotry4041(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(gomock.Any()).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/-1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReposiotry4042(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(gomock.Any()).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/bad/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(uint(1), "/child_dir/test.txt").Return("../utils/test_files/file_tree/child_dir/test.txt", nil)

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/file/child_dir/test.txt", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "attachment; filename=test.txt", responseRecorder.Result().Header.Get("Content-Disposition"))

	assert.Equal(t, "text/plain", responseRecorder.Result().Header.Get("Content-Type"))
}

func TestGetFileFromRepository4041(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/bad/file/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository4042(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/-1/file/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository4045(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(uint(1), "/bad").Return("", errors.New("err"))

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/file/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository500(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(uint(1), "/bad").Return("nonexistent", nil)

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/file/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetTreeFromRepository200(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(uint(1)).Return(fileTree, nil, nil)

	req, _ := http.NewRequest("GET", "/api/v2/versions/1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

	receivedTree := make(map[string]int64)
	responseJSON, _ := io.ReadAll(responseRecorder.Body)
	_ = json.Unmarshal(responseJSON, &receivedTree)

	assert.Equal(t, fileTree, receivedTree)
}

func TestGetTreeFromRepository4001(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(gomock.Any()).Return(fileTree, nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/bad/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetTreeFromRepository4002(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(gomock.Any()).Return(fileTree, nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v2/versions/-1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

// CreateMultipartFile bundles a file into an object go can interact with to return it as a response.
// Returns file, content-type, and error
func CreateMultipartFile(filePath string) (io.Reader, string, error) {
	// create a buffer to hold the file in memory
	body := new(bytes.Buffer)

	mwriter := multipart.NewWriter(body)
	defer mwriter.Close()

	w, err := mwriter.CreateFormFile("file", filePath)

	if err != nil {
		return body, "", err
	}

	in, err := os.Open(filePath)

	if err != nil {
		return body, "", err
	}

	defer in.Close()

	_, err = io.Copy(w, in)

	if err != nil {
		return body, "", err
	}

	return body, mwriter.FormDataContentType(), nil
}
