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

	mockBranchRepository              *mocks.MockModelRepositoryInterface[*models.Branch]
	mockProjectPostRepository         *mocks.MockModelRepositoryInterface[*models.ProjectPost]
	mockReviewRepository              *mocks.MockModelRepositoryInterface[*models.Review]
	mockBranchCollaboratorRepository  *mocks.MockModelRepositoryInterface[*models.BranchCollaborator]
	mockDiscussionContainerRepository *mocks.MockModelRepositoryInterface[*models.DiscussionContainer]
	mockDiscussionRepository          *mocks.MockModelRepositoryInterface[*models.Discussion]
	mockFilesystem                    *mocks.MockFilesystem

	pendingBranch *models.Branch
	successBranch *models.Branch
	failedBranch  *models.Branch

	projectPost *models.ProjectPost

	c   *gin.Context
	cwd string
)

func TestMain(m *testing.M) {
	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
