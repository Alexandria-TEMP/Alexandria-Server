package services

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
)

func beforeEachRender(t *testing.T) {
	t.Helper()

	// setup models
	pendingBranch = &models.Branch{RenderStatus: models.Pending}
	successBranch = &models.Branch{RenderStatus: models.Success}
	failedBranch = &models.Branch{RenderStatus: models.Failure}
	projectPost = &models.ProjectPost{}

	// Setup mock DB and vfs
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBranchRepository = mocks.NewMockModelRepositoryInterface[*models.Branch](mockCtrl)
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	mockBranchService = mocks.NewMockBranchService(mockCtrl)

	// Create render service
	renderService = RenderService{
		BranchRepository:      mockBranchRepository,
		ProjectPostRepository: mockProjectPostRepository,
		Filesystem:            mockFilesystem,
		BranchService:         mockBranchService,
	}
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwd, "render"))
}

func TestRenderSuccess1(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testRenderSuccessTemplate(t, "good_quarto_project_1")
}

func TestRenderSuccess2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testRenderSuccessTemplate(t, "good_quarto_project_2")
}

func TestRenderSuccess3(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testRenderSuccessTemplate(t, "good_quarto_project_3")
}

func testRenderSuccessTemplate(t *testing.T, dirName string) {
	t.Helper()
	beforeEachRender(t)
	defer cleanup(t)

	pendingBranch.ID = 10
	successBranch.ID = 10
	dirPath := filepath.Join(cwd, "..", "utils", "test_files", dirName)
	renderDirPath := filepath.Join(cwd, "render")

	mockBranchRepository.EXPECT().Update(pendingBranch).Return(pendingBranch, nil)
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(renderDirPath).AnyTimes()
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(true, "").Times(1)
	mockFilesystem.EXPECT().CreateCommit().Return(nil).Times(1)
	mockBranchRepository.EXPECT().Update(successBranch).Return(successBranch, nil).Times(1)

	renderService.RenderBranch(pendingBranch)

	_, err := os.Stat(renderDirPath)
	assert.Nil(t, err)
}

func TestRenderUnzipFailed(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	pendingBranch.ID = 10
	failedBranch.ID = 10
	renderDirPath := filepath.Join(cwd, "render")

	mockBranchRepository.EXPECT().Update(pendingBranch).Return(pendingBranch, nil)
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().Unzip().Return(errors.New("failed")).Times(1)
	mockBranchRepository.EXPECT().Update(failedBranch).Return(failedBranch, nil).Times(1)
	mockFilesystem.EXPECT().Reset()

	renderService.RenderBranch(pendingBranch)

	_, err := os.Stat(renderDirPath)
	assert.NotNil(t, err)
}

func TestRenderExistsFailed(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	beforeEachRender(t)
	defer cleanup(t)

	pendingBranch.ID = 10
	failedBranch.ID = 10
	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "good_quarto_project_1")
	renderDirPath := filepath.Join(cwd, "render")

	mockBranchRepository.EXPECT().Update(pendingBranch).Return(pendingBranch, nil)
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).AnyTimes()
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(renderDirPath).AnyTimes()
	mockFilesystem.EXPECT().Unzip().Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(false, "").Times(1)
	mockBranchRepository.EXPECT().Update(failedBranch).Return(failedBranch, nil).Times(1)

	renderService.RenderBranch(pendingBranch)
}

func TestIsValidProjectNoYamlorYml(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "bad_quarto_project_1")

	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).Times(2)

	assert.False(t, renderService.IsValidProject())
}

func TestIsValidProjectNotDefaultType(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "bad_quarto_project_4")

	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).Times(2)

	assert.False(t, renderService.IsValidProject())
}

func TestIsValidProjectWithYaml(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "good_quarto_project_4")

	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).Times(2)

	assert.True(t, renderService.IsValidProject())
}

func TestGetRenderFileSuccess(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	projectPostID := uint(99)
	renderFilePath := filepath.Join(cwd, "..", "utils", "test_files", "good_repository_setup", "render", "test.html")
	successBranch.ID = 0
	successBranch.ProjectPostID = &projectPostID
	projectPost.ID = 99
	projectPost.PostID = 100

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(successBranch, nil).Times(1)
	mockBranchService.EXPECT().GetBranchProjectPost(successBranch).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(99)).Return(projectPost, nil).Times(1)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(100)).Times(1)
	mockFilesystem.EXPECT().CheckoutBranch("0").Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(true, "").Times(1)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(renderFilePath)

	returnedPath, err202, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.Nil(t, err404)
	assert.Equal(t, renderFilePath, returnedPath)
}

func TestGetRenderFileNoBranch(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(successBranch, errors.New("failed")).Times(1)

	_, err202, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}

func TestGetRenderFileNoProjectPost(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	projectPostID := uint(99)
	successBranch.ID = 0
	successBranch.ProjectPostID = &projectPostID

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(successBranch, nil).Times(1)
	mockBranchService.EXPECT().GetBranchProjectPost(successBranch).Return(projectPost, errors.New("failed"))

	_, err202, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}

func TestGetRenderFilePending(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	projectPostID := uint(99)
	pendingBranch.ID = 0
	pendingBranch.ProjectPostID = &projectPostID

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(pendingBranch, nil).Times(1)
	mockBranchService.EXPECT().GetBranchProjectPost(pendingBranch).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(99)).Return(projectPost, nil).Times(1)

	_, err202, err404 := renderService.GetRenderFile(pendingBranch.ID)

	assert.NotNil(t, err202)
	assert.Nil(t, err404)
}

func TestGetRenderFileFailed(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	projectPostID := uint(99)
	failedBranch.ID = 0
	failedBranch.ProjectPostID = &projectPostID

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(failedBranch, nil).Times(1)
	mockBranchService.EXPECT().GetBranchProjectPost(failedBranch).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(99)).Return(projectPost, nil).Times(1)

	_, err202, err404 := renderService.GetRenderFile(failedBranch.ID)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}

func TestGetRenderNoGitBranch(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	projectPostID := uint(99)
	successBranch.ID = 0
	successBranch.ProjectPostID = &projectPostID
	projectPost.ID = 99
	projectPost.PostID = 100

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(successBranch, nil).Times(1)
	mockBranchService.EXPECT().GetBranchProjectPost(successBranch).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(99)).Return(projectPost, nil).Times(1)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(100)).Times(1)
	mockFilesystem.EXPECT().CheckoutBranch("0").Return(errors.New("failed")).Times(1)

	_, err202, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}

func TestGetRenderDoesntExist(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	projectPostID := uint(99)
	successBranch.ID = 0
	successBranch.ProjectPostID = &projectPostID
	projectPost.ID = 99
	projectPost.PostID = 100

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(successBranch, nil).Times(1)
	mockBranchService.EXPECT().GetBranchProjectPost(successBranch).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(99)).Return(projectPost, nil).Times(1)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(100)).Times(1)
	mockFilesystem.EXPECT().CheckoutBranch("0").Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return(false, "").Times(1)
	mockBranchRepository.EXPECT().Update(successBranch).Return(successBranch, nil)

	_, err202, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.NotNil(t, err404)
}
