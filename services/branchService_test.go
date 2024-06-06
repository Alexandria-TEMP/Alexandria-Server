package services

import (
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
)

func beforeEachBranch(t *testing.T) {
	t.Helper()

	// setup models

	// Setup mock DB and vfs
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBranchRepository = mocks.NewMockModelRepositoryInterface[*models.Branch](mockCtrl)
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockReviewRepository = mocks.NewMockModelRepositoryInterface[*models.Review](mockCtrl)
	mockBranchCollaboratorRepository = mocks.NewMockModelRepositoryInterface[*models.BranchCollaborator](mockCtrl)
	mockDiscussionContainerRepository = mocks.NewMockModelRepositoryInterface[*models.DiscussionContainer](mockCtrl)
	mockDiscussionRepository = mocks.NewMockModelRepositoryInterface[*models.Discussion](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)

	// Create branch service
	branchService = BranchService{
		BranchRepository:              mockBranchRepository,
		ProjectPostRepository:         mockProjectPostRepository,
		ReviewRepository:              mockReviewRepository,
		BranchCollaboratorRepository:  mockBranchCollaboratorRepository,
		DiscussionContainerRepository: mockDiscussionContainerRepository,
		DiscussionRepository:          mockDiscussionRepository,
		Filesystem:                    mockFilesystem,
	}
}
