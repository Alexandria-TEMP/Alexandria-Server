package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	renderService RenderService
	branchService BranchService

	mockRenderService             *mocks.MockRenderService
	mockPostCollaboratorService   *mocks.MockPostCollaboratorService
	mockBranchCollaboratorService *mocks.MockBranchCollaboratorService
	mockBranchService             *mocks.MockBranchService

	mockBranchRepository              *mocks.MockModelRepositoryInterface[*models.Branch]
	mockClosedBranchRepository        *mocks.MockModelRepositoryInterface[*models.ClosedBranch]
	mockPostRepository                *mocks.MockModelRepositoryInterface[*models.Post]
	mockProjectPostRepository         *mocks.MockModelRepositoryInterface[*models.ProjectPost]
	mockReviewRepository              *mocks.MockModelRepositoryInterface[*models.BranchReview]
	mockBranchCollaboratorRepository  *mocks.MockModelRepositoryInterface[*models.BranchCollaborator]
	mockPostCollaboratorRepository    *mocks.MockModelRepositoryInterface[*models.PostCollaborator]
	mockDiscussionContainerRepository *mocks.MockModelRepositoryInterface[*models.DiscussionContainer]
	mockDiscussionRepository          *mocks.MockModelRepositoryInterface[*models.Discussion]
	mockMemberRepository              *mocks.MockModelRepositoryInterface[*models.Member]

	mockFilesystem *mocks.MockFilesystem

	pendingBranch *models.Branch
	successBranch *models.Branch
	failedBranch  *models.Branch

	memberA, memberB, memberC models.Member

	projectPost *models.ProjectPost

	c   *gin.Context
	cwd string
)

func setupTestSuite() {
}

func teardownTestSuite() {
}

func TestMain(m *testing.M) {
	setupTestSuite()

	cwd, _ = os.Getwd()
	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	code := m.Run()

	teardownTestSuite()
	os.Exit(code)
}
