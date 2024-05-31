package services

import (
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
)

func beforeEach(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	versionService = VersionService{Filesystem: mockFilesystem}
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwd, "render"))
}

func TestCreateVersionSuccess1(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testGoodProjectTemplate(t, "good_quarto_project_1")
}

func TestCreateVersionSuccess2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testGoodProjectTemplate(t, "good_quarto_project_2")
}

func TestCreateVersionSuccess3(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testGoodProjectTemplate(t, "good_quarto_project_3")
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

// func TestGetRenderFileSuccess(t *testing.T) {
// }

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

	assert.NotNil(t, err)

	renderDirPath := filepath.Join(cwd, "render")
	_, err = os.Stat(renderDirPath)
	assert.NotNil(t, err)
}

func testGoodProjectTemplate(t *testing.T, dirName string) {
	t.Helper()
	beforeEach(t)
	defer cleanup(t)

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(uint(0), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().CountRenderFiles().Return(1).Times(1)
	mockFilesystem.EXPECT().RemoveProjectDirectory().Return(nil).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)
	assert.Equal(t, &exampleVersion, version)

	// Wait until model has completed rendering
	for version.RenderStatus == models.RenderPending {
		print()
	}
	assert.Equal(t, models.RenderSuccess, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render")
	_, err = os.Stat(renderDirPath)
	assert.Nil(t, err)
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
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)
	assert.Equal(t, &exampleVersion, version)

	// Wait until model has completed rendering
	for version.RenderStatus == models.RenderPending {
		print()
	}
	assert.Equal(t, models.RenderFailure, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render", "quarto_project.html")
	assert.Equal(t, false, filesystem.FileExists(renderDirPath))
}
