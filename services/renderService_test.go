package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func beforeEachRender(t *testing.T) {
	t.Helper()

	// setup models
	_ = lock.Lock()
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
	mockPostRepository = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)

	// Create render service
	renderService = RenderService{
		BranchRepository:      mockBranchRepository,
		PostRepository:        mockPostRepository,
		ProjectPostRepository: mockProjectPostRepository,
		Filesystem:            mockFilesystem,
		BranchService:         mockBranchService,
	}
}

func cleanup(t *testing.T) {
	t.Helper()

	_ = lock.Unlock()

	_ = os.RemoveAll(filepath.Join(cwd, "render"))
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
	mockFilesystem.EXPECT().RenderExists().Return("", nil).Times(1)
	mockFilesystem.EXPECT().CreateCommit().Return(nil).Times(1)
	mockBranchRepository.EXPECT().Update(successBranch).Return(successBranch, nil).Times(1)

	renderService.RenderBranch(pendingBranch, lock, mockFilesystem)

	_, err := os.Stat(renderDirPath)
	assert.Nil(t, err)
	assert.False(t, lock.Locked())
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

	renderService.RenderBranch(pendingBranch, lock, mockFilesystem)

	_, err := os.Stat(renderDirPath)
	assert.NotNil(t, err)
	assert.False(t, lock.Locked())
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
	mockFilesystem.EXPECT().RenderExists().Return("", fmt.Errorf("oh no")).Times(1)
	mockBranchRepository.EXPECT().Update(failedBranch).Return(failedBranch, nil).Times(1)

	renderService.RenderBranch(pendingBranch, lock, mockFilesystem)
	assert.False(t, lock.Locked())
}

func TestIsValidProjectNoYamlorYml(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "bad_quarto_project_1")

	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).Times(2)

	assert.False(t, renderService.isValidProject(mockFilesystem))
}

func TestIsValidProjectNotDefaultType(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "bad_quarto_project_4")

	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).Times(2)

	assert.False(t, renderService.isValidProject(mockFilesystem))
}

func TestIsValidProjectWithYaml(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	dirPath := filepath.Join(cwd, "..", "utils", "test_files", "good_quarto_project_4")

	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(dirPath).Times(2)

	assert.True(t, renderService.isValidProject(mockFilesystem))
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
	mockFilesystem.EXPECT().LockDirectory(projectPost.PostID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(100)).Times(1)
	mockFilesystem.EXPECT().CheckoutBranch("0").Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return("", nil).Times(1)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return(renderFilePath)

	returnedPath, outputLock, err202, err204, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.Nil(t, err204)
	assert.Nil(t, err404)
	assert.Equal(t, renderFilePath, returnedPath)
	assert.Equal(t, lock, outputLock)
	assert.True(t, lock.Locked())
}

func TestGetRenderFileNoBranch(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	mockBranchRepository.EXPECT().GetByID(uint(0)).Return(successBranch, errors.New("failed")).Times(1)

	_, _, err202, err204, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.Nil(t, err204)
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

	_, _, err202, err204, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.Nil(t, err204)
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

	_, _, err202, err204, err404 := renderService.GetRenderFile(pendingBranch.ID)

	assert.NotNil(t, err202)
	assert.Nil(t, err204)
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

	_, _, err202, err204, err404 := renderService.GetRenderFile(failedBranch.ID)

	assert.Nil(t, err202)
	assert.NotNil(t, err204)
	assert.Nil(t, err404)
}

func TestGetRenderFileNoGitBranch(t *testing.T) {
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
	mockFilesystem.EXPECT().LockDirectory(projectPost.PostID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(100)).Times(1)
	mockFilesystem.EXPECT().CheckoutBranch("0").Return(errors.New("failed")).Times(1)

	_, _, err202, err204, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.Nil(t, err204)
	assert.NotNil(t, err404)
	assert.False(t, lock.Locked())
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
	mockFilesystem.EXPECT().LockDirectory(projectPost.PostID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(100)).Times(1)
	mockFilesystem.EXPECT().CheckoutBranch("0").Return(nil).Times(1)
	mockFilesystem.EXPECT().RenderExists().Return("", fmt.Errorf("oh no")).Times(1)
	mockBranchRepository.EXPECT().Update(successBranch).Return(successBranch, nil)

	_, _, err202, err204, err404 := renderService.GetRenderFile(successBranch.ID)

	assert.Nil(t, err202)
	assert.Nil(t, err204)
	assert.NotNil(t, err404)
	assert.False(t, lock.Locked())
}

func TestGetMainRenderFileGoodWeather(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	// Setup data
	postID := uint(10)

	post := &models.Post{
		Model:        gorm.Model{ID: postID},
		RenderStatus: models.Success,
	}

	// Setup mocks
	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil)
	mockFilesystem.EXPECT().LockDirectory(uint(10)).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)

	// Checking out master branch will succeed
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)

	// The render will always exist
	mockFilesystem.EXPECT().RenderExists().Return("render_filename", nil)
	mockFilesystem.EXPECT().GetCurrentRenderDirPath().Return("path")

	// Function under test
	actualPath, outputLock, err202, err204, err404 := renderService.GetMainRenderFile(postID)
	if err202 != nil || err204 != nil || err404 != nil {
		t.Fatal(err202, err204, err404)
	}

	expectedPath := "path/render_filename"

	assert.Equal(t, expectedPath, actualPath)
	assert.Equal(t, lock, outputLock)
	assert.True(t, lock.Locked())
}

func TestGetMainRenderStillPending(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	// Setup data
	postID := uint(10)

	post := &models.Post{
		Model:        gorm.Model{ID: postID},
		RenderStatus: models.Pending,
	}

	// Setup mocks
	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil)

	// Function under test
	_, _, err1, err2, err3 := renderService.GetMainRenderFile(postID)

	assert.NotNil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
}

func TestGetMainRenderFailed(t *testing.T) {
	beforeEachRender(t)
	defer cleanup(t)

	// Setup data
	postID := uint(10)

	post := &models.Post{
		Model:        gorm.Model{ID: postID},
		RenderStatus: models.Failure,
	}

	// Setup mocks
	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil)

	// Function under test
	_, _, err1, err2, err3 := renderService.GetMainRenderFile(postID)

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.Nil(t, err3)
}
