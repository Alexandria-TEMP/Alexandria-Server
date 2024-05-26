package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
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

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), uint(1)).Return(&examplePendingVersion, nil).Times(1)

	zipPath := filepath.Join(cwd, "..", "utils", "test_files", "file_handling_test.zip")
	body, dataType, err := filesystem.CreateMultipartFile(zipPath)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/version/1", body)
	req.Header.Add("Content-Type", dataType)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCreateVersion4001(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), uint(1)).Return(&examplePendingVersion, nil).Times(0)

	req, _ := http.NewRequest("POST", "/api/v1/version/1", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateVersion4002(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), uint(1)).Return(&examplePendingVersion, nil).Times(0)

	zipPath := filepath.Join(cwd, "..", "utils", "test_files", "file_handling_test.zip")
	body, dataType, err := filesystem.CreateMultipartFile(zipPath)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/version/bad", body)
	req.Header.Add("Content-Type", dataType)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestCreateVersion500(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), uint(1)).Return(&examplePendingVersion, errors.New("err")).Times(1)

	zipPath := filepath.Join(cwd, "..", "utils", "test_files", "file_handling_test.zip")
	body, dataType, err := filesystem.CreateMultipartFile(zipPath)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/api/v1/version/1", body)
	req.Header.Add("Content-Type", dataType)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetRender200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(uint(1), uint(0)).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/render", http.NoBody)
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

	mockVersionService.EXPECT().GetRenderFile(uint(1), uint(0)).Return("", errors.New("err"), nil)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusAccepted, responseRecorder.Result().StatusCode)
}

func TestGetRender4001(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/bad/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender4002(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/bad/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender4003(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/-1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender4004(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/render/1234.html", nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/-1/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetRender404(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRenderFile(uint(1), uint(0)).Return("", nil, errors.New("err"))

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/render", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetReposiotry200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(uint(1), uint(0)).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/repository", http.NoBody)
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

	mockVersionService.EXPECT().GetRepositoryFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/-1/1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReposiotry4042(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/bad/1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReposiotry4043(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/5/-1/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetReposiotry4044(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetRepositoryFile(gomock.Any(), gomock.Any()).Return("../utils/test_files/good_repository_setup/quarto_project.zip", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/5/bad/repository", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository200(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(uint(1), uint(0), "/child-dir/_child.qmd").Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/blob/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, "attachment; filename=_child.qmd", responseRecorder.Result().Header.Get("Content-Disposition"))

	fileHeader := make([]byte, headerSize)
	_, err := responseRecorder.Result().Body.Read(fileHeader)
	assert.Nil(t, err)

	fileContentType := http.DetectContentType(fileHeader)

	assert.Equal(t, fileContentType, "application/octet-stream")
}

func TestGetFileFromRepository4041(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(gomock.Any(), gomock.Any(), gomock.Any()).Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/bad/1/blob/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository4042(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(gomock.Any(), gomock.Any(), gomock.Any()).Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/-1/1/blob/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository4043(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(gomock.Any(), gomock.Any(), gomock.Any()).Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/bad/blob/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository4044(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(gomock.Any(), gomock.Any(), gomock.Any()).Return("../utils/test_files/good_quarto_project_4/child-dir/_child.qmd", nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/-1/blob/child-dir/_child.qmd", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository4045(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(uint(1), uint(0), "/bad").Return("", errors.New("err"))

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/blob/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode)
}

func TestGetFileFromRepository500(t *testing.T) {
	beforeEachVersion(t)

	mockVersionService.EXPECT().GetFileFromRepository(uint(1), uint(0), "/bad").Return("nonexistent", nil)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/blob/bad", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode)
}

func TestGetTreeFromRepository200(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(uint(1), uint(0)).Return(fileTree, nil, nil)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/1/tree", http.NoBody)
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
	mockVersionService.EXPECT().GetTreeFromRepository(gomock.Any(), gomock.Any()).Return(fileTree, nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/bad/1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetTreeFromRepository4002(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(gomock.Any(), gomock.Any()).Return(fileTree, nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/-1/1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetTreeFromRepository4003(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(gomock.Any(), gomock.Any()).Return(fileTree, nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/bad/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}

func TestGetTreeFromRepository4004(t *testing.T) {
	beforeEachVersion(t)

	fileTree := map[string]int64{"file1": 4}
	mockVersionService.EXPECT().GetTreeFromRepository(gomock.Any(), gomock.Any()).Return(fileTree, nil, nil).Times(0)

	req, _ := http.NewRequest("GET", "/api/v1/version/0/-1/tree", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	defer responseRecorder.Result().Body.Close()

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode)
}
