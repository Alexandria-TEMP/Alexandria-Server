package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

// Variables that are used by all service tests
// IMPORTANT: tests are responsible for initializing their values!

// Mocked repositories
var (
	c                                 *gin.Context
	cwd                               string
	mockPostRepository                *mocks.MockModelRepositoryInterface[*models.Post]
	mockProjectPostRepository         *mocks.MockModelRepositoryInterface[*models.ProjectPost]
	mockMemberRepository              *mocks.MockModelRepositoryInterface[*models.Member]
	mockPostCollaboratorRepository    *mocks.MockModelRepositoryInterface[*models.PostCollaborator]
	mockBranchCollaboratorRepository  *mocks.MockModelRepositoryInterface[*models.BranchCollaborator]
	mockDiscussionRepository          *mocks.MockModelRepositoryInterface[*models.Discussion]
	mockDiscussionContainerRepository *mocks.MockModelRepositoryInterface[*models.DiscussionContainer]
)

// Mocked services
var (
	mockPostCollaboratorService   *mocks.MockPostCollaboratorService
	mockBranchCollaboratorService *mocks.MockBranchCollaboratorService
)

// Data that can be used by tests
var (
	memberA, memberB, memberC models.Member
	discussionA               models.Discussion
	discussionContainerA      models.DiscussionContainer
	memberService             MemberService
	exampleMember             models.Member
	exampleMemberDTO          models.MemberDTO
	exampleSTag1              *tags.ScientificFieldTag
	exampleSTag2              *tags.ScientificFieldTag
)

func setupTestSuite() {
}

func teardownTestSuite() {
}

func TestMain(m *testing.M) {
	tag1 := tags.ScientificFieldTag{
		ScientificField: "Mathematics",
		Subtags:         []*tags.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag1 = &tag1
	tag2 := tags.ScientificFieldTag{
		ScientificField: "",
		Subtags:         []*tags.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag2 = &tag2
	scientificFieldTagArray := []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2}
	exampleMember = models.Member{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: scientificFieldTagArray,
		},
	}
	exampleMemberDTO = models.MemberDTO{
		FirstName:             "John",
		LastName:              "Smith",
		Email:                 "john.smith@gmail.com",
		Password:              "password",
		Institution:           "TU Delft",
		ScientificFieldTagIDs: []uint{},
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	setupTestSuite()

	code := m.Run()

	teardownTestSuite()

	os.Exit(code)
}
