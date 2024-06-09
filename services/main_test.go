package services

import (
	"os"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

// Variables that are used by all service tests
// Tests are responsible for initializing their values

var (
	postRepositoryMock               *mocks.MockModelRepositoryInterface[*models.Post]
	projectPostRepositoryMock        *mocks.MockModelRepositoryInterface[*models.ProjectPost]
	memberRepositoryMock             *mocks.MockModelRepositoryInterface[*models.Member]
	postCollaboratorRepositoryMock   *mocks.MockModelRepositoryInterface[*models.PostCollaborator]
	branchCollaboratorRepositoryMock *mocks.MockModelRepositoryInterface[*models.BranchCollaborator]
)

var (
	postCollaboratorServiceMock   *mocks.MockPostCollaboratorService
	branchCollaboratorServiceMock *mocks.MockBranchCollaboratorService
)

var memberA, memberB, memberC models.Member

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
