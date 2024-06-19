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
	memberService MemberService

	mockRenderService             *mocks.MockRenderService
	mockPostCollaboratorService   *mocks.MockPostCollaboratorService
	mockBranchCollaboratorService *mocks.MockBranchCollaboratorService
	mockBranchService             *mocks.MockBranchService
	mockTagService                *mocks.MockTagService

	mockBranchRepository                      *mocks.MockModelRepositoryInterface[*models.Branch]
	mockClosedBranchRepository                *mocks.MockModelRepositoryInterface[*models.ClosedBranch]
	mockPostRepository                        *mocks.MockModelRepositoryInterface[*models.Post]
	mockProjectPostRepository                 *mocks.MockModelRepositoryInterface[*models.ProjectPost]
	mockBranchReviewRepository                *mocks.MockModelRepositoryInterface[*models.BranchReview]
	mockBranchCollaboratorRepository          *mocks.MockModelRepositoryInterface[*models.BranchCollaborator]
	mockPostCollaboratorRepository            *mocks.MockModelRepositoryInterface[*models.PostCollaborator]
	mockDiscussionContainerRepository         *mocks.MockModelRepositoryInterface[*models.DiscussionContainer]
	mockDiscussionRepository                  *mocks.MockModelRepositoryInterface[*models.Discussion]
	mockMemberRepository                      *mocks.MockModelRepositoryInterface[*models.Member]
	mockScientificFieldTagRepository          *mocks.MockModelRepositoryInterface[*models.ScientificFieldTag]
	mockScientificFieldTagContainerRepository *mocks.MockModelRepositoryInterface[*models.ScientificFieldTagContainer]

	mockFilesystem *mocks.MockFilesystem

	pendingBranch *models.Branch
	successBranch *models.Branch
	failedBranch  *models.Branch

	memberA, memberB, memberC models.Member

	discussionA          models.Discussion
	discussionContainerA models.DiscussionContainer

	exampleMember    models.Member
	exampleMemberDTO models.MemberDTO
	exampleSTag1     *models.ScientificFieldTag
	exampleSTag2     *models.ScientificFieldTag

	projectPost *models.ProjectPost

	c   *gin.Context
	cwd string
)

func setupTestSuite() {
}

func teardownTestSuite() {
}

func TestMain(m *testing.M) {
	tag1 := models.ScientificFieldTag{
		ScientificField: "Mathematics",
		Subtags:         []*models.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag1 = &tag1
	tag2 := models.ScientificFieldTag{
		ScientificField: "",
		Subtags:         []*models.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag2 = &tag2
	scientificFieldTagArray := []*models.ScientificFieldTag{exampleSTag1, exampleSTag2}
	exampleMember = models.Member{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
		ScientificFieldTagContainer: models.ScientificFieldTagContainer{
			ScientificFieldTags: scientificFieldTagArray,
		},
	}
	exampleMemberDTO = models.MemberDTO{
		FirstName:                     "John",
		LastName:                      "Smith",
		Email:                         "john.smith@gmail.com",
		Password:                      "password",
		Institution:                   "TU Delft",
		ScientificFieldTagContainerID: 0,
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	setupTestSuite()

	cwd, _ = os.Getwd()
	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	code := m.Run()

	teardownTestSuite()
	os.Exit(code)
}
