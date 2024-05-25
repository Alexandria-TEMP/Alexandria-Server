package services

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
)

func beforeEach(t *testing.T) {
	t.Helper()

	// Create fresh repo
	var err error
	db, err = database.InitializeTestDatabase()
	if err != nil {
		log.Fatalf("Could not initialize test database: %s", err)
	}

	versionRepository = database.ModelRepository[*models.Version]{Database: db}

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	versionService = VersionService{
		VersionRepository: versionRepository,
		Filesystem:        mockFilesystem,
	}
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwd, "render"))

	db.Unscoped().Where("id >= 0").Delete(&models.Post{})
	db.Unscoped().Where("id >= 0").Delete(&models.Version{})
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

// Can take a while, so if this times out increase limit
func TestCreateVersionSuccess4(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

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

func TestGetRenderFileSuccess(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	versionService.VersionRepository.Create(&successVersion)

	mockFilesystem.EXPECT().SetCurrentVersion(uint(2), uint(0)).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(true, "").Times(1)
	mockFilesystem.EXPECT().GetRenderFile().Return([]byte{53, 54, 55, 56}, nil).Times(1)

	file, err202, err404 := versionService.GetRender(2, 0)

	assert.Nil(t, err202)
	assert.Nil(t, err404)
	assert.Equal(t, []byte{53, 54, 55, 56}, file)
}

func TestGetRenderFileFailure1(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	versionService.VersionRepository.Create(&pendingVersion)

	mockFilesystem.EXPECT().SetCurrentVersion(uint(0), uint(0)).Times(0)
	mockFilesystem.EXPECT().RenderExists().Times(0)
	mockFilesystem.EXPECT().GetRenderFile().Times(0)

	_, err202, err404 := versionService.GetRender(0, 0)

	assert.NotNil(t, err202)
	assert.Nil(t, err404)
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
	mockFilesystem.EXPECT().RenderExists().Return(true, "").Times(1)
	mockFilesystem.EXPECT().RemoveProjectDirectory().Return(nil).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)
	assert.Equal(t, &pendingVersion, version)

	// Wait until model has completed rendering
	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Success, version.RenderStatus)

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
	assert.Equal(t, &pendingVersion, version)

	// Wait until model has completed rendering
	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Failure, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render", "quarto_project.html")
	assert.Equal(t, false, filesystem.FileExists(renderDirPath))
}
