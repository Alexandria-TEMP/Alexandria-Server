package services

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func beforeEachBranch(t *testing.T) {
	t.Helper()

	// setup models
	pendingBranch = &models.Branch{RenderStatus: models.Pending}
	successBranch = &models.Branch{RenderStatus: models.Success}
	failedBranch = &models.Branch{RenderStatus: models.Failure}
	projectPost = &models.ProjectPost{}

	// Setup mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRenderService = mocks.NewMockRenderService(mockCtrl)
	mockBranchRepository = mocks.NewMockModelRepositoryInterface[*models.Branch](mockCtrl)
	mockClosedBranchRepository = mocks.NewMockModelRepositoryInterface[*models.ClosedBranch](mockCtrl)
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockReviewRepository = mocks.NewMockModelRepositoryInterface[*models.BranchReview](mockCtrl)
	mockBranchCollaboratorRepository = mocks.NewMockModelRepositoryInterface[*models.BranchCollaborator](mockCtrl)
	mockDiscussionContainerRepository = mocks.NewMockModelRepositoryInterface[*models.DiscussionContainer](mockCtrl)
	mockDiscussionRepository = mocks.NewMockModelRepositoryInterface[*models.Discussion](mockCtrl)
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	mockBranchCollaboratorService = mocks.NewMockBranchCollaboratorService(mockCtrl)

	// Create branch service
	branchService = BranchService{
		RenderService:                 mockRenderService,
		BranchRepository:              mockBranchRepository,
		ProjectPostRepository:         mockProjectPostRepository,
		ReviewRepository:              mockReviewRepository,
		BranchCollaboratorRepository:  mockBranchCollaboratorRepository,
		DiscussionContainerRepository: mockDiscussionContainerRepository,
		DiscussionRepository:          mockDiscussionRepository,
		MemberRepository:              mockMemberRepository,
		Filesystem:                    mockFilesystem,
		BranchCollaboratorService:     mockBranchCollaboratorService,
		ClosedBranchRepository:        mockClosedBranchRepository,
	}
}

func TestGetBranchSuccess(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(9)).Return(successBranch, nil)

	branch, err := branchService.GetBranch(uint(9))
	assert.Nil(t, err)
	assert.Equal(t, *successBranch, branch)
}

func TestGetBranchFailed(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(9)).Return(successBranch, errors.New("failed"))

	_, err := branchService.GetBranch(uint(9))
	assert.NotNil(t, err)
}

func TestCreateBranchSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	projectPost.PostID = 12
	collaborator := &models.BranchCollaborator{MemberID: 12}
	expectedBranch := &models.Branch{
		Collaborators:             []*models.BranchCollaborator{collaborator},
		ProjectPostID:             10,
		RenderStatus:              models.Success,
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}
	outputBranch := &models.Branch{
		Collaborators:             []*models.BranchCollaborator{collaborator},
		ProjectPostID:             10,
		RenderStatus:              models.Success,
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}
	newProjectPost := &models.ProjectPost{
		Model:        gorm.Model{ID: 10},
		PostID:       12,
		OpenBranches: []*models.Branch{expectedBranch},
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().Create(&models.DiscussionContainer{}).Return(nil)
	mockProjectPostRepository.EXPECT().Update(newProjectPost).Return(newProjectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(12))
	mockFilesystem.EXPECT().CreateBranch("0")
	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(12)).Return(collaborator, nil)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{12}, false).Return([]*models.BranchCollaborator{collaborator}, nil)

	branch, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs: []uint{12},
		ProjectPostID:          10,
	})

	assert.Nil(t, err404)
	assert.Nil(t, err500)
	assert.Equal(t, outputBranch, &branch)
}

func TestCreateBranchNoProjectPost(t *testing.T) {
	beforeEachBranch(t)

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, errors.New("failed"))

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs: []uint{12, 11},
		ProjectPostID:          10,
	})

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

func TestCreateBranchFailedUpdateProjectPost(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	projectPost.PostID = 12
	expectedBranch := &models.Branch{
		Collaborators:             []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID:             10,
		RenderStatus:              models.Success,
		BranchOverallReviewStatus: models.BranchOpenForReview,
		DiscussionContainer:       models.DiscussionContainer{},
		UpdatedScientificFields:   []models.ScientificField{},
		Reviews:                   []*models.BranchReview{},
	}
	newProjectPost := &models.ProjectPost{
		Model:        gorm.Model{ID: 10},
		PostID:       12,
		OpenBranches: []*models.Branch{expectedBranch},
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().Create(&models.DiscussionContainer{}).Return(nil)
	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Return(newProjectPost, errors.New("failed"))
	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(12)).Return(&models.BranchCollaborator{MemberID: 12}, nil)
	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(11)).Return(&models.BranchCollaborator{MemberID: 11}, nil)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{12, 11}, false).Return([]*models.BranchCollaborator{{MemberID: 12, BranchID: 11}}, nil)

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs: []uint{12, 11},
		ProjectPostID:          10,
	})

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestCreateBranchFailedGit(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	projectPost.PostID = 12
	expectedBranch := &models.Branch{
		Collaborators:             []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID:             10,
		RenderStatus:              models.Success,
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}
	newProjectPost := &models.ProjectPost{
		Model:        gorm.Model{ID: 10},
		PostID:       12,
		OpenBranches: []*models.Branch{expectedBranch},
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().Create(&models.DiscussionContainer{}).Return(nil)
	mockBranchRepository.EXPECT().Create(expectedBranch).DoAndReturn(
		func(expectedBranch *models.Branch) error {
			expectedBranch.ID = 15
			return nil
		})
	mockFilesystem.EXPECT().CheckoutDirectory(uint(12))
	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Return(newProjectPost, nil)
	mockFilesystem.EXPECT().CreateBranch("0").Return(errors.New("failed"))
	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(12)).Return(&models.BranchCollaborator{MemberID: 12}, nil)
	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(11)).Return(&models.BranchCollaborator{MemberID: 12}, nil)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{12, 11}, false).Return([]*models.BranchCollaborator{{MemberID: 12, BranchID: 11}}, nil)

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs: []uint{12, 11},
		ProjectPostID:          10,
	})

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestDeleteBranchSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().DeleteBranch("10").Return(nil)
	mockBranchRepository.EXPECT().Delete(uint(10)).Return(nil)

	assert.Nil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedDeleteGitBranch(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().DeleteBranch("10").Return(errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedDelete(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().DeleteBranch("10").Return(nil)
	mockBranchRepository.EXPECT().Delete(uint(10)).Return(errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestGetReviewStatusSuccess(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:            gorm.Model{ID: 10},
		UpdatedPostTitle: "title",
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{{BranchReviewDecision: models.Approved}, {BranchReviewDecision: models.Rejected}}, nil)

	decisions, err := branchService.GetReviewStatus(uint(10))
	assert.Nil(t, err)
	assert.Equal(t, []models.BranchReviewDecision{models.Approved, models.Rejected}, decisions)
}

func TestGetReviewStatusFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:            gorm.Model{ID: 10},
		UpdatedPostTitle: "title",
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	_, err := branchService.GetReviewStatus(uint(10))
	assert.NotNil(t, err)
}

func TestCreateReviewSuccess(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		ReviewingMemberID:    11,
		BranchReviewDecision: models.Approved,
	}
	expected := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Approved,
	}
	branch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}
	newBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		Reviews:                   []*models.BranchReview{expected},
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockReviewRepository.EXPECT().Create(expected).Return(nil)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, nil)

	branchreview, err := branchService.CreateReview(form)
	assert.Nil(t, err)
	assert.Equal(t, expected, &branchreview)
}

func TestCreateReviewSuccessMerge(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		ReviewingMemberID:    11,
		BranchReviewDecision: models.Approved,
	}
	expected := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Approved,
	}
	branch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		BranchOverallReviewStatus: models.BranchOpenForReview,
		ProjectPostID:             5,
	}
	newBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		Reviews:                   []*models.BranchReview{expected, expected, expected},
		BranchOverallReviewStatus: models.BranchPeerReviewed,
		ProjectPostID:             5,
	}
	closed := &models.ClosedBranch{
		Branch:               *newBranch,
		SupercededBranch:     &models.Branch{Model: gorm.Model{ID: 50}},
		ProjectPostID:        5,
		BranchReviewDecision: models.Approved,
	}
	projectPost.ID = 5
	projectPost.OpenBranches = append(projectPost.OpenBranches, branch)
	projectPost.LastMergedBranch = &models.Branch{Model: gorm.Model{ID: 50}}
	newProjectPost := &models.ProjectPost{
		Model:            gorm.Model{ID: 5},
		LastMergedBranch: newBranch,
		OpenBranches:     []*models.Branch{},
		ClosedBranches:   []*models.ClosedBranch{closed},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{expected, expected}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockFilesystem.EXPECT().Merge("10", "master").Return(nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Return(newProjectPost, nil)

	branchreview, err := branchService.CreateReview(form)
	assert.Nil(t, err)
	assert.Equal(t, expected, &branchreview)
	assert.Equal(t, models.BranchPeerReviewed, branch.BranchOverallReviewStatus)
	assert.Equal(t, projectPost.LastMergedBranch, newBranch)
}

func TestCreateReviewSuccessReject(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		ReviewingMemberID:    11,
		BranchReviewDecision: models.Rejected,
	}
	approval := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Approved,
	}
	expected := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Rejected,
	}
	branch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		BranchOverallReviewStatus: models.BranchOpenForReview,
		ProjectPostID:             5,
	}
	newBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		Reviews:                   []*models.BranchReview{approval, approval, expected},
		BranchOverallReviewStatus: models.BranchPeerReviewed,
		ProjectPostID:             5,
	}
	closed := &models.ClosedBranch{
		Branch:               *newBranch,
		ProjectPostID:        5,
		BranchReviewDecision: models.Rejected,
	}
	projectPost.ID = 5
	projectPost.OpenBranches = append(projectPost.OpenBranches, branch)
	projectPost.LastMergedBranch = &models.Branch{Model: gorm.Model{ID: 50}}
	newProjectPost := &models.ProjectPost{
		Model:            gorm.Model{ID: 5},
		OpenBranches:     []*models.Branch{},
		ClosedBranches:   []*models.ClosedBranch{closed},
		LastMergedBranch: &models.Branch{Model: gorm.Model{ID: 50}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{approval, approval}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Return(newProjectPost, nil)

	branchreview, err := branchService.CreateReview(form)
	assert.Nil(t, err)
	assert.Equal(t, expected, &branchreview)
	assert.Equal(t, models.BranchRejected, branch.BranchOverallReviewStatus)
	assert.Equal(t, &models.Branch{Model: gorm.Model{ID: 50}}, newProjectPost.LastMergedBranch)
}

func TestCreateReviewFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		ReviewingMemberID:    11,
		BranchReviewDecision: models.Approved,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	_, err := branchService.CreateReview(form)
	assert.NotNil(t, err)
}

func TestCreateReviewFailedGetMember(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		ReviewingMemberID:    11,
		BranchReviewDecision: models.Approved,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, errors.New("failed"))

	_, err := branchService.CreateReview(form)
	assert.NotNil(t, err)
}

func TestCreateReviewFailedUpdateBranch(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		ReviewingMemberID:    11,
		BranchReviewDecision: models.Approved,
	}
	expected := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Approved,
	}
	branch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}
	newBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		Reviews:                   []*models.BranchReview{expected},
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{}, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockReviewRepository.EXPECT().Create(expected).Return(nil)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, errors.New("failed"))

	_, err := branchService.CreateReview(form)
	assert.NotNil(t, err)
}

func TestMemberCanReviewSuccessTrue(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model:            gorm.Model{ID: 11},
		ScientificFields: []models.ScientificField{models.Mathematics},
	}
	branch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		ProjectPostID:             20,
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}
	projectPost := &models.ProjectPost{
		Model: gorm.Model{ID: 20},
		Post:  models.Post{ScientificFields: []models.ScientificField{models.Mathematics, models.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockReviewRepository.EXPECT().GetBy(&models.BranchReview{BranchID: 10}).Return([]*models.BranchReview{}, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)

	canReview, err := branchService.MemberCanReview(10, 11)
	assert.Nil(t, err)
	assert.True(t, canReview)
}

func TestMemberCanReviewSuccessFalse(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model:            gorm.Model{ID: 11},
		ScientificFields: []models.ScientificField{models.Mathematics},
	}
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 20,
	}
	projectPost := &models.ProjectPost{
		Model: gorm.Model{ID: 20},
		Post:  models.Post{ScientificFields: []models.ScientificField{models.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)

	canReview, err := branchService.MemberCanReview(10, 11)
	assert.Nil(t, err)
	assert.False(t, canReview)
}

func TestMemberCanReviewFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	_, err := branchService.MemberCanReview(10, 11)
	assert.NotNil(t, err)
}

func TestMemberCanReviewFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 20,
	}
	projectPost := &models.ProjectPost{
		Model: gorm.Model{ID: 20},
		Post:  models.Post{ScientificFields: []models.ScientificField{models.Mathematics, models.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, errors.New("failed"))

	_, err := branchService.MemberCanReview(10, 11)
	assert.NotNil(t, err)
}

func TestMemberCanReviewFailedGetMember(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 20,
	}
	projectPost := &models.ProjectPost{
		Model: gorm.Model{ID: 20},
		Post:  models.Post{ScientificFields: []models.ScientificField{models.Mathematics, models.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, errors.New("failed"))

	_, err := branchService.MemberCanReview(10, 11)
	assert.NotNil(t, err)
}

func TestGetProjectSuccess(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	expectedFilePath := "../utils/test_files/good_repository_setup/quarto_project.zip"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentZipFilePath().Return(expectedFilePath)

	filePath, err := branchService.GetProject(10)
	assert.Nil(t, err)
	assert.Equal(t, expectedFilePath, filePath)
}

func TestGetProjectFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(nil, errors.New("failed"))

	filePath, err := branchService.GetProject(10)
	assert.NotNil(t, err)
	assert.Equal(t, "", filePath)
}

func TestGetProjectFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(nil, errors.New("failed"))

	filePath, err := branchService.GetProject(10)
	assert.NotNil(t, err)
	assert.Equal(t, "", filePath)
}

func TestGetProjectFailedCheckoutBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))

	filePath, err := branchService.GetProject(10)
	assert.NotNil(t, err)
	assert.Equal(t, "", filePath)
}

func TestUploadProjectSuccess(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		RenderStatus:  models.Success,
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	expected := &models.Branch{
		RenderStatus:  models.Pending,
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	file := &multipart.FileHeader{}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(nil)
	mockFilesystem.EXPECT().SaveZipFile(gomock.Any(), file).Return(nil)
	mockFilesystem.EXPECT().CreateCommit().Return(nil)
	mockBranchRepository.EXPECT().Update(expected).Return(expected, nil)
	mockRenderService.EXPECT().RenderBranch(branch)

	assert.Nil(t, branchService.UploadProject(c, file, 10))
}

func TestUploadProjectFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(nil, errors.New("failed"))

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(nil, errors.New("failed"))

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedCheckoutBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedCleanDir(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(errors.New("failed"))

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedSaveZipFile(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
		RenderStatus:  models.Success,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	file := &multipart.FileHeader{}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(nil)
	mockFilesystem.EXPECT().SaveZipFile(gomock.Any(), file).Return(errors.New("failed"))
	mockBranchRepository.EXPECT().Update(gomock.Any()).Return(branch, nil)
	mockFilesystem.EXPECT().Reset()

	err := branchService.UploadProject(c, file, 10)
	assert.NotNil(t, err)
}

func TestGetFiletreeSuccess(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	expectedFileTree := map[string]int64{
		"file1.txt": 1234,
		"file2.txt": 5678,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetFileTree().Return(expectedFileTree, nil)

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, expectedFileTree, fileTree)
}

func TestGetFiletreeFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(nil, errors.New("failed"))

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.NotNil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetFiletreeFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(nil, errors.New("failed"))

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.NotNil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetFiletreeFailedCheckoutBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.NotNil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetFiletreeFailedGetFileTree(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetFileTree().Return(nil, errors.New("failed"))

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetFileFromProjectSuccess(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	relFilepath := "/child_dir/test.txt"
	quartoDirPath := "../utils/test_files/file_tree"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(quartoDirPath)

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(quartoDirPath, relFilepath), absFilepath)
}

func TestGetFileFromProjectRelativePathContainsDotDot(t *testing.T) {
	beforeEachBranch(t)

	relFilepath := "../some/unsafe/path"

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}

func TestGetFileFromProjectFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	relFilepath := "example.qmd"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(nil, errors.New("failed"))

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}

func TestGetFileFromProjectFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	relFilepath := "example.qmd"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(nil, errors.New("failed"))

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}

func TestGetFileFromProjectFailedCheckoutBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	relFilepath := "child_dir/test.txt"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}

func TestGetFileFromProjectFileDoesNotExist(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 5,
	}
	projectPost := &models.ProjectPost{
		PostID: 50,
	}
	relFilepath := "nonexistent/file.txt"
	quartoDirPath := "../utils/test_files/good_repository_setup"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(quartoDirPath)

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}
