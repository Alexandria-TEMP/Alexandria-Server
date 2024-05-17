package services_tests

import (
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"go.uber.org/mock/gomock"
)

func beforeEach(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	versionService = services.VersionService{Filesystem: mockFilesystem}
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwd, "render"))
}

// func TestSaveRepository200(t *testing.T) {
// 	beforeEach(t)

// 	body, dataType := filesystem.CreateMultipartFile("file.zip")

// 	req, _ := http.NewRequest("POST", "/api/v1/version/1", body)
// 	req.Header.Add("Content-Type", dataType)
// 	router.ServeHTTP(responseRecorder, req)

// 	defer responseRecorder.Result().Body.Close()

// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

// 	// cleanup(t)
// }

// func TestRenderSuccess(t *testing.T) {
// 	beforeEach(t)

// 	versionController.VersionService.GetFilesystem().SetCurrentVersion(0, 0)
// 	err := versionController.VersionService.RenderProject()

// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
// }

func TestCreateVersionSuccess1(t *testing.T) {
	testGoodProjectTemplate(t, "good_quarto_project_1")
}

func TestCreateVersionSuccess2(t *testing.T) {
	testGoodProjectTemplate(t, "good_quarto_project_2")
}

func TestCreateVersionSuccess3(t *testing.T) {
	testGoodProjectTemplate(t, "good_quarto_project_3")
}

// Can take a while, so if this times out increase limit
func TestCreateVersionSuccess4(t *testing.T) {
	testGoodProjectTemplate(t, "good_quarto_project_4")
}

func TestCreateVersionDelayedFailure1(t *testing.T) {
	testBadProjectTemplate(t, "bad_quarto_project_1")
}

func TestCreateVersionDelayedFailure2(t *testing.T) {
	testBadProjectTemplate(t, "bad_quarto_project_2")
}

func TestCreateVersionDelayedFailure3(t *testing.T) {
	testBadProjectTemplate(t, "bad_quarto_project_3")
}

func TestCreateVersionDelayedFailure4(t *testing.T) {
	testBadProjectTemplate(t, "bad_quarto_project_4")
}

func TestCreateVersionImmediateFailure(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(uint(0), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(errors.New("")).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(0)
	mockFilesystem.EXPECT().RemoveProjectDirectory().Return(nil).Times(0)
	mockFilesystem.EXPECT().RemoveRepository().Return(nil).Times(1)

	_, err := versionService.CreateVersion(c, file, 2)

	assert.NotEqual(t, nil, err)

	renderDirPath := filepath.Join(cwd, "render")
	_, err = os.Stat(renderDirPath)
	assert.NotEqual(t, nil, err)
}

func testGoodProjectTemplate(t *testing.T, dirName string) {
	t.Helper()
	beforeEach(t)
	defer cleanup(t)

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(uint(0), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().RemoveProjectDirectory().Return(nil).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "util", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Equal(t, nil, err)
	assert.Equal(t, exampleVersion, version)

	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Success, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render")
	_, err = os.Stat(renderDirPath)
	assert.Equal(t, nil, err)
}

func testBadProjectTemplate(t *testing.T, dirName string) {
	t.Helper()
	beforeEach(t)
	defer cleanup(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(uint(0), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().RemoveRepository().Return(nil).Times(1)
	mockFilesystem.EXPECT().RemoveProjectDirectory().Return(nil).Times(0)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "util", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Equal(t, nil, err)
	assert.Equal(t, exampleVersion, version)

	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Failure, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render", "quarto_project.html")
	assert.Equal(t, false, services.FileExists(renderDirPath))
}
