package services

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
	"go.uber.org/mock/gomock"
)

func beforeEach(t *testing.T) {
	t.Helper()

	// Setup mock DB
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockVersionRepository = mocks.NewMockRepositoryInterface[*models.Version](mockCtrl)

	// Setup mock filesystem
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)

	// Cretae version service
	versionService = VersionService{
		VersionRepository: mockVersionRepository,
		Filesystem:        mockFilesystem,
	}
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

// Can take a while, so if this times out increase limit
// func TestCreateVersionSuccess4(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip()
// 	}

// 	testGoodProjectTemplate(t, "good_quarto_project_4")
// }

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

func TestCreateVersionDelayedFailure5(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(gomock.Any(), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(errors.New("err")).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", "bad_quarto_project_1")).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()
	mockFilesystem.EXPECT().RenderExists().Times(0)
	mockFilesystem.EXPECT().RemoveRepository().Times(1)

	mockVersionRepository.EXPECT().Create(&models.Version{RenderStatus: models.Pending})
	mockVersionRepository.EXPECT().Update(&models.Version{RenderStatus: models.Failure}).Times(1)

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)

	assert.Equal(t, models.Pending, version.RenderStatus)

	// Wait until model has completed rendering
	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Failure, version.RenderStatus)
}

func TestCreateVersionDelayedFailure6(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(gomock.Any(), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", "good_quarto_project_1")).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()
	mockFilesystem.EXPECT().RenderExists().Return(false, "").Times(1)
	mockFilesystem.EXPECT().RemoveRepository().Times(0)

	mockVersionRepository.EXPECT().Create(&models.Version{RenderStatus: models.Pending})
	mockVersionRepository.EXPECT().Update(&models.Version{RenderStatus: models.Failure}).Times(1)

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)

	assert.Equal(t, models.Pending, version.RenderStatus)

	// Wait until model has completed rendering
	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Failure, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render", "quarto_project.html")
	assert.Equal(t, false, utils.FileExists(renderDirPath))
}

func TestCreateVersionImmediateFailure(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(gomock.Any(), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(errors.New("")).Times(1)
	mockFilesystem.EXPECT().RemoveRepository().Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(0)

	mockVersionRepository.EXPECT().Create(&models.Version{RenderStatus: models.Pending})
	mockVersionRepository.EXPECT().Update(&models.Version{RenderStatus: models.Failure}).Times(1)

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

	mockFilesystem.EXPECT().SetCurrentVersion(gomock.Any(), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(true, "").Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()
	mockFilesystem.EXPECT().RemoveRepository().Times(0)

	mockVersionRepository.EXPECT().Create(&models.Version{RenderStatus: models.Pending}).Times(1)
	mockVersionRepository.EXPECT().Update(&models.Version{RenderStatus: models.Success}).Times(1)

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)
	assert.Equal(t, models.Pending, version.RenderStatus)

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

	file := &multipart.FileHeader{}

	mockFilesystem.EXPECT().SetCurrentVersion(gomock.Any(), uint(2)).Times(1)
	mockFilesystem.EXPECT().SaveRepository(c, file).Return(nil).Times(1)
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", dirName)).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(filepath.Join(cwd, "render")).AnyTimes()
	mockFilesystem.EXPECT().RenderExists().Times(0)
	mockFilesystem.EXPECT().RemoveRepository().Times(1)

	mockVersionRepository.EXPECT().Create(&models.Version{RenderStatus: models.Pending})
	mockVersionRepository.EXPECT().Update(&models.Version{RenderStatus: models.Failure}).Times(1)

	version, err := versionService.CreateVersion(c, file, 2)

	assert.Nil(t, err)

	assert.Equal(t, models.Pending, version.RenderStatus)

	// Wait until model has completed rendering
	for version.RenderStatus == models.Pending {
		print()
	}
	assert.Equal(t, models.Failure, version.RenderStatus)

	renderDirPath := filepath.Join(cwd, "render", "quarto_project.html")
	assert.Equal(t, false, utils.FileExists(renderDirPath))
}

func TestGetRenderFileSuccess(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(true, "").Times(1)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return("test").Times(1)

	mockVersionRepository.EXPECT().GetByID(successVersion.ID).Return(&successVersion, nil).Times(1)
	mockVersionRepository.EXPECT().Update(gomock.Any()).Times(0)

	path, err202, err404 := versionService.GetRenderFile(successVersion.ID, 0)

	assert.Nil(t, err202)
	assert.Nil(t, err404)
	assert.Equal(t, "test", path)
}

func TestGetRenderFileFailure1(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(pendingVersion.ID, uint(0)).Times(0)
	mockFilesystem.EXPECT().RenderExists().Times(0)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return("").Times(0)

	mockVersionRepository.EXPECT().GetByID(pendingVersion.ID).Return(&pendingVersion, nil).Times(1)
	mockVersionRepository.EXPECT().Update(gomock.Any()).Times(0)

	_, err202, err404 := versionService.GetRenderFile(pendingVersion.ID, 0)

	assert.NotNil(t, err202)
	assert.Nil(t, err404)
}

func TestGetRenderFileFailure2(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(failureVersion.ID, uint(0)).Times(0)
	mockFilesystem.EXPECT().RenderExists().Times(0)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return("").Times(0)

	mockVersionRepository.EXPECT().GetByID(failureVersion.ID).Return(&failureVersion, nil).Times(1)
	mockVersionRepository.EXPECT().Update(gomock.Any()).Times(0)

	_, err202, err404 := versionService.GetRenderFile(failureVersion.ID, 0)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}

func TestGetRenderFileFailure3(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(false, "").Times(1)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return("test").Times(0)

	mockVersionRepository.EXPECT().GetByID(successVersion.ID).Return(&successVersion, nil).Times(1)
	mockVersionRepository.EXPECT().Update(&successVersion).Times(1)

	_, err202, err404 := versionService.GetRenderFile(successVersion.ID, 0)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}

func TestGetTreeFromRepositorySuccess(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return("../utils/test_files/file_tree").Times(1)
	mockFilesystem.EXPECT().GetFileTree().Return(map[string]int64{"child_dir/test.txt": 0, "example.qmd": 0}, nil).Times(1)

	tree, err404, err500 := versionService.GetTreeFromRepository(successVersion.ID, 0)

	assert.Nil(t, err404)
	assert.Nil(t, err500)
	assert.Equal(t, map[string]int64{"child_dir/test.txt": 0, "example.qmd": 0}, tree)
}

func TestGetTreeFromRepositoryFailure1(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return("doesntexist").Times(1)
	mockFilesystem.EXPECT().GetFileTree().Return(map[string]int64{"child_dir/test.txt": 0, "example.qmd": 0}, nil).Times(0)

	_, err404, err500 := versionService.GetTreeFromRepository(successVersion.ID, 0)

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

func TestGetTreeFromRepositoryFailure2(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return("../utils/test_files/file_tree").Times(1)
	mockFilesystem.EXPECT().GetFileTree().Return(nil, errors.New("err")).Times(1)

	_, err404, err500 := versionService.GetTreeFromRepository(successVersion.ID, 0)

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestGetRepositoryFileSuccess(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentZipFilePath().Return("../utils/test_files/good_repository_setup/quarto_project.zip").Times(2)

	path, err := versionService.GetRepositoryFile(successVersion.ID, 0)

	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(cwd, "..", "utils", "test_files", "good_repository_setup", "quarto_project.zip"), path)
}

func TestGetRepositoryFileFailure(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentZipFilePath().Return("doesntexist").Times(1)

	_, err := versionService.GetRepositoryFile(successVersion.ID, 0)

	assert.NotNil(t, err)
}

func TestGetFileFromRepositorySuccess(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", "file_tree")).Times(1)

	absFilepath, err := versionService.GetFileFromRepository(successVersion.ID, 0, "/child_dir/test.txt")

	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(cwd, "..", "utils", "test_files", "file_tree", "child_dir", "test.txt"), absFilepath)
}

func TestGetFileFromRepositoryFailure(t *testing.T) {
	beforeEach(t)
	defer cleanup(t)

	mockFilesystem.EXPECT().SetCurrentVersion(successVersion.ID, uint(0)).Times(1)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(filepath.Join(cwd, "..", "utils", "test_files", "file_tree")).Times(1)

	_, err := versionService.GetFileFromRepository(successVersion.ID, 0, "../../../.env")

	assert.NotNil(t, err)
}
