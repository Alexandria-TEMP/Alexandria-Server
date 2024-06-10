package services

import (
	"os"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

// Variables that are used by all service tests
// IMPORTANT: tests are responsible for initializing their values!

// Mocked repositories
var (
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
)

func setupTestSuite() {
}

func teardownTestSuite() {
}

func TestMain(m *testing.M) {
	setupTestSuite()

	code := m.Run()

	teardownTestSuite()

	os.Exit(code)
}
