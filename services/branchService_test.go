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
	mockPostRepository = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockBranchReviewRepository = mocks.NewMockModelRepositoryInterface[*models.BranchReview](mockCtrl)
	mockDiscussionContainerRepository = mocks.NewMockModelRepositoryInterface[*models.DiscussionContainer](mockCtrl)
	mockDiscussionRepository = mocks.NewMockModelRepositoryInterface[*models.Discussion](mockCtrl)
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	mockBranchCollaboratorService = mocks.NewMockBranchCollaboratorService(mockCtrl)
	mockTagService = mocks.NewMockTagService(mockCtrl)
	mockPostCollaboratorService = mocks.NewMockPostCollaboratorService(mockCtrl)

	// Create branch service
	branchService = BranchService{
		BranchRepository:              mockBranchRepository,
		PostRepository:                mockPostRepository,
		ProjectPostRepository:         mockProjectPostRepository,
		ReviewRepository:              mockBranchReviewRepository,
		DiscussionContainerRepository: mockDiscussionContainerRepository,
		DiscussionRepository:          mockDiscussionRepository,
		MemberRepository:              mockMemberRepository,
		Filesystem:                    mockFilesystem,
		ClosedBranchRepository:        mockClosedBranchRepository,
		BranchCollaboratorService:     mockBranchCollaboratorService,
		PostCollaboratorService:       mockPostCollaboratorService,
		RenderService:                 mockRenderService,
		TagService:                    mockTagService,
	}
}

func TestGetBranchSuccess(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(9)).Return(successBranch, nil)

	branch, err := branchService.GetBranch(uint(9))
	assert.Nil(t, err)
	assert.Equal(t, successBranch, branch)
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
	outputBranch := &models.Branch{
		Collaborators:                      []*models.BranchCollaborator{collaborator},
		ProjectPostID:                      &projectPost.ID,
		RenderStatus:                       models.Success,
		BranchOverallReviewStatus:          models.BranchOpenForReview,
		UpdatedScientificFieldTagContainer: &models.ScientificFieldTagContainer{ScientificFieldTags: []*models.ScientificFieldTag{}},
	}
	newProjectPost := &models.ProjectPost{
		Model:        gorm.Model{ID: 10},
		PostID:       12,
		OpenBranches: []*models.Branch{outputBranch},
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().Create(&models.DiscussionContainer{}).Return(nil)
	mockProjectPostRepository.EXPECT().Update(newProjectPost).Return(newProjectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(12))
	mockFilesystem.EXPECT().CreateBranch("0")
	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(12)).Return(collaborator, nil)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{12}, false).Return([]*models.BranchCollaborator{collaborator}, nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)

	branch, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs:    []uint{12},
		ProjectPostID:             10,
		UpdatedScientificFieldIDs: []uint{},
	}, &models.Member{Model: gorm.Model{ID: 12}})

	assert.Nil(t, err404)
	assert.Nil(t, err500)
	assert.Equal(t, outputBranch, branch)
}

func TestCreateBranchNoProjectPost(t *testing.T) {
	beforeEachBranch(t)

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, errors.New("failed"))

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs: []uint{12, 11},
		ProjectPostID:          10,
	}, &models.Member{Model: gorm.Model{ID: 12}})

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

func TestCreateBranchFailedUpdateProjectPost(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	projectPost.PostID = 12
	expectedBranch := &models.Branch{
		Collaborators:                      []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID:                      &projectPost.ID,
		RenderStatus:                       models.Success,
		BranchOverallReviewStatus:          models.BranchOpenForReview,
		DiscussionContainer:                models.DiscussionContainer{},
		UpdatedScientificFieldTagContainer: &models.ScientificFieldTagContainer{},
		Reviews:                            []*models.BranchReview{},
	}
	newProjectPost := &models.ProjectPost{
		Model:        gorm.Model{ID: 10},
		PostID:       12,
		OpenBranches: []*models.Branch{expectedBranch},
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().Create(&models.DiscussionContainer{}).Return(nil)
	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Return(newProjectPost, errors.New("failed"))
	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(12)).Return(&models.BranchCollaborator{MemberID: 12}, nil)
	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(11)).Return(&models.BranchCollaborator{MemberID: 11}, nil)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{12, 11}, false).Return([]*models.BranchCollaborator{{MemberID: 12, BranchID: 11}}, nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs:    []uint{12, 11},
		ProjectPostID:             10,
		UpdatedScientificFieldIDs: []uint{},
	}, &models.Member{Model: gorm.Model{ID: 12}})

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestCreateBranchFailedGit(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	projectPost.PostID = 12
	expectedBranch := &models.Branch{
		Collaborators:                      []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID:                      &projectPost.ID,
		RenderStatus:                       models.Success,
		BranchOverallReviewStatus:          models.BranchOpenForReview,
		UpdatedScientificFieldTagContainer: &models.ScientificFieldTagContainer{ScientificFieldTags: []*models.ScientificFieldTag{}},
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
	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(12)).Return(&models.BranchCollaborator{MemberID: 12}, nil)
	mockBranchCollaboratorService.EXPECT().GetBranchCollaborator(uint(11)).Return(&models.BranchCollaborator{MemberID: 12}, nil)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{12, 11}, false).Return([]*models.BranchCollaborator{{MemberID: 12, BranchID: 11}}, nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		CollaboratingMemberIDs:    []uint{12, 11},
		ProjectPostID:             10,
		UpdatedScientificFieldIDs: []uint{},
	}, &models.Member{Model: gorm.Model{ID: 12}})

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestDeleteBranchSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50
	projectPost.Post = models.Post{DiscussionContainerID: 1}
	projectPost.ID = 5
	projectPostID := uint(5)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(1)).Return(&models.DiscussionContainer{Model: gorm.Model{ID: 1}}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 5})
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().DeleteBranch("10").Return(nil)
	mockBranchRepository.EXPECT().Delete(uint(10)).Return(nil)

	assert.Nil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50
	projectPostID := uint(5)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50
	projectPostID := uint(5)

	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedDeleteGitBranch(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50
	projectPost.Post = models.Post{DiscussionContainerID: 1}
	projectPost.ID = 5
	projectPostID := uint(5)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(1)).Return(&models.DiscussionContainer{Model: gorm.Model{ID: 1}}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 5})
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().DeleteBranch("10").Return(errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestDeleteBranchFailedDelete(t *testing.T) {
	beforeEachBranch(t)

	projectPost.PostID = 50
	projectPost.Post = models.Post{DiscussionContainerID: 1}
	projectPost.ID = 5
	projectPostID := uint(5)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(1)).Return(&models.DiscussionContainer{Model: gorm.Model{ID: 1}}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 5})
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().DeleteBranch("10").Return(nil)
	mockBranchRepository.EXPECT().Delete(uint(10)).Return(errors.New("failed"))

	assert.NotNil(t, branchService.DeleteBranch(10))
}

func TestGetReviewStatusSuccess(t *testing.T) {
	beforeEachBranch(t)

	approved := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Approved,
	}
	rejected := &models.BranchReview{
		// Model:          gorm.Model{ID: 1},
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 12}},
		BranchReviewDecision: models.Rejected,
	}
	branch := &models.Branch{
		Model:   gorm.Model{ID: 10},
		Reviews: []*models.BranchReview{approved, rejected},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	decisions, err := branchService.GetAllBranchReviewStatuses(uint(10))
	assert.Nil(t, err)
	assert.Equal(t, []models.BranchReviewDecision{models.Approved, models.Rejected}, decisions)
}

func TestGetReviewStatusFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	_, err := branchService.GetAllBranchReviewStatuses(uint(10))
	assert.NotNil(t, err)
}

func TestCreateReviewSuccess(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
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
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockBranchReviewRepository.EXPECT().Create(expected).Return(nil)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, nil)

	branchreview, err := branchService.CreateReview(form, &models.Member{Model: gorm.Model{ID: 11}})
	assert.Nil(t, err)
	assert.Equal(t, expected, branchreview)
}

func TestCreateReviewSuccessMergeDoesntSupercede(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(5)
	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
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
		ProjectPostID:             &projectPostID,
		Reviews:                   []*models.BranchReview{expected, expected},
	}
	newBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		Reviews:                   []*models.BranchReview{expected, expected, expected},
		BranchOverallReviewStatus: models.BranchPeerReviewed,
		ProjectPostID:             nil,
	}
	closed := &models.ClosedBranch{
		Branch:               *newBranch,
		SupercededBranch:     &models.Branch{Model: gorm.Model{ID: 50}},
		ProjectPostID:        5,
		BranchReviewDecision: models.Approved,
	}
	projectPost.ID = 5
	projectPost.OpenBranches = append(projectPost.OpenBranches, branch)
	newProjectPost := &models.ProjectPost{
		Model:          gorm.Model{ID: 5},
		OpenBranches:   []*models.Branch{},
		ClosedBranches: []*models.ClosedBranch{closed},
	}
	discussions := &models.DiscussionContainer{Model: gorm.Model{ID: 1}}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockBranchReviewRepository.EXPECT().Create(expected).Return(nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(0))
	mockFilesystem.EXPECT().Merge("10", "master").Return(nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(projectPost, nil)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(discussions, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 5})
	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Return(newProjectPost, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{
		ProjectPostID:        5,
		BranchReviewDecision: models.Approved,
	}).Return(nil, nil)
	mockPostCollaboratorService.EXPECT().MergeContributors(projectPost, newBranch.Collaborators)
	mockPostCollaboratorService.EXPECT().MergeReviewers(projectPost, newBranch.Reviews)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, nil)
	mockPostRepository.EXPECT().Update(&models.Post{DiscussionContainer: *discussions}).Return(&models.Post{DiscussionContainer: *discussions}, nil)

	branchreview, err := branchService.CreateReview(form, &models.Member{Model: gorm.Model{ID: 11}})
	assert.Nil(t, err)
	assert.Equal(t, expected, branchreview)
	assert.Equal(t, models.BranchPeerReviewed, branch.BranchOverallReviewStatus)
	assert.Nil(t, projectPost.ClosedBranches[0].SupercededBranchID)
}

func TestCreateReviewSuccessMergeSupercedes(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(8)
	form := forms.ReviewCreationForm{
		BranchID:             10,
		BranchReviewDecision: models.Approved,
	}
	expectedReview := &models.BranchReview{
		BranchID:             10,
		Member:               models.Member{Model: gorm.Model{ID: 11}},
		BranchReviewDecision: models.Approved,
	}
	oldApprovedBranch := &models.ClosedBranch{
		Model:         gorm.Model{ID: 11},
		ProjectPostID: 8,
		Branch: models.Branch{
			Model:                     gorm.Model{ID: 9},
			BranchOverallReviewStatus: models.BranchPeerReviewed,
			Reviews:                   []*models.BranchReview{expectedReview, expectedReview},
		},
	}
	initialBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		BranchOverallReviewStatus: models.BranchOpenForReview,
		Reviews:                   []*models.BranchReview{expectedReview, expectedReview},
		ProjectPostID:             &projectPostID,
	}
	expectedBranch := &models.Branch{
		Model:                     gorm.Model{ID: 10},
		BranchOverallReviewStatus: models.BranchPeerReviewed,
		Reviews:                   []*models.BranchReview{expectedReview, expectedReview, expectedReview},
		ProjectPostID:             nil,
	}
	expectedClosedBranch := &models.ClosedBranch{
		Branch:               *expectedBranch,
		SupercededBranch:     &oldApprovedBranch.Branch,
		ProjectPostID:        5,
		BranchReviewDecision: models.Approved,
	}
	initialProjectPost := &models.ProjectPost{
		Model:          gorm.Model{ID: 5},
		OpenBranches:   []*models.Branch{initialBranch},
		ClosedBranches: []*models.ClosedBranch{oldApprovedBranch},
		PostID:         7,
	}
	expectedProjectPost := &models.ProjectPost{
		Model:            gorm.Model{ID: 5},
		OpenBranches:     []*models.Branch{},
		ClosedBranches:   []*models.ClosedBranch{oldApprovedBranch, expectedClosedBranch},
		PostID:           7,
		PostReviewStatus: models.Reviewed,
	}
	discussions := &models.DiscussionContainer{}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(initialBranch, nil)
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(&models.Member{Model: gorm.Model{ID: 11}}, nil)
	mockBranchReviewRepository.EXPECT().Create(expectedReview).Return(nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(8)).Return(initialProjectPost, nil)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(discussions, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 5}).Return([]*models.ClosedBranch{oldApprovedBranch}, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(7))
	mockFilesystem.EXPECT().Merge("10", "master")
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{
		ProjectPostID:        5,
		BranchReviewDecision: models.Approved,
	}).Return([]*models.ClosedBranch{oldApprovedBranch}, nil)
	mockBranchRepository.EXPECT().Update(expectedBranch).Return(expectedBranch, nil)
	mockPostCollaboratorService.EXPECT().MergeContributors(initialProjectPost, expectedBranch.Collaborators)
	mockPostCollaboratorService.EXPECT().MergeReviewers(initialProjectPost, expectedBranch.Reviews)
	mockProjectPostRepository.EXPECT().Update(expectedProjectPost).Return(expectedProjectPost, nil)
	mockPostRepository.EXPECT().Update(&models.Post{DiscussionContainer: *discussions}).Return(&models.Post{DiscussionContainer: *discussions}, nil)

	_, _ = branchService.CreateReview(form, &models.Member{Model: gorm.Model{ID: 11}})

	assert.Equal(t, expectedProjectPost, initialProjectPost)
}

func TestCreateReviewFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
		BranchReviewDecision: models.Approved,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, errors.New("failed"))

	_, err := branchService.CreateReview(form, &models.Member{Model: gorm.Model{ID: 11}})
	assert.NotNil(t, err)
}

func TestCreateReviewFailedUpdateBranch(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:             10,
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
	mockMemberRepository.EXPECT().GetByID(uint(11)).Return(member, nil)
	mockBranchReviewRepository.EXPECT().Create(expected).Return(nil)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, errors.New("failed"))

	_, err := branchService.CreateReview(form, &models.Member{Model: gorm.Model{ID: 11}})
	assert.NotNil(t, err)
}

func TestGetProjectSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	expectedFilePath := "../utils/test_files/good_repository_setup/quarto_project.zip"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentZipFilePath().Return(expectedFilePath)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

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

	projectPostID := uint(5)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(nil, errors.New("failed"))

	filePath, err := branchService.GetProject(10)
	assert.NotNil(t, err)
	assert.Equal(t, "", filePath)
}

func TestGetProjectFailedCheckoutBranch(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	filePath, err := branchService.GetProject(10)
	assert.NotNil(t, err)
	assert.Equal(t, "", filePath)
}

func TestUploadProjectSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		RenderStatus:  models.Success,
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	expected := &models.Branch{
		RenderStatus:  models.Pending,
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	file := &multipart.FileHeader{}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(nil)
	mockFilesystem.EXPECT().SaveZipFile(gomock.Any(), file).Return(nil)
	mockFilesystem.EXPECT().CreateCommit().Return(nil)
	mockBranchRepository.EXPECT().Update(expected).Return(expected, nil)
	mockRenderService.EXPECT().RenderBranch(branch)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	assert.Nil(t, branchService.UploadProject(c, file, 10))
}

func TestUploadProjectFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(nil, errors.New("failed"))

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedGetProjectPost(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(5)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(5)).Return(nil, errors.New("failed"))

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedCheckoutBranch(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedCleanDir(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(errors.New("failed"))
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	assert.NotNil(t, branchService.UploadProject(c, nil, 10))
}

func TestUploadProjectFailedSaveZipFile(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
		RenderStatus:  models.Success,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	file := &multipart.FileHeader{}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(nil)
	mockFilesystem.EXPECT().SaveZipFile(gomock.Any(), file).Return(errors.New("failed"))
	mockBranchRepository.EXPECT().Update(gomock.Any()).Return(branch, nil)
	mockFilesystem.EXPECT().Reset()
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	err := branchService.UploadProject(c, file, 10)
	assert.NotNil(t, err)
}

func TestGetFiletreeSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	expectedFileTree := map[string]int64{
		"file1.txt": 1234,
		"file2.txt": 5678,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetFileTree().Return(expectedFileTree, nil)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

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

	projectPostID := uint(5)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
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

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.NotNil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetFiletreeFailedGetFileTree(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetFileTree().Return(nil, errors.New("failed"))
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	fileTree, err1, err2 := branchService.GetFiletree(10)
	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetFileFromProjectSuccess(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	relFilepath := "/child_dir/test.txt"
	quartoDirPath := "../utils/test_files/file_tree"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(quartoDirPath)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

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

	projectPostID := uint(5)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
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

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	relFilepath := "child_dir/test.txt"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(errors.New("failed"))
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}

func TestGetFileFromProjectFileDoesNotExist(t *testing.T) {
	beforeEachBranch(t)

	projectPostID := uint(20)
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: &projectPostID,
	}
	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: 20},
		PostID: 50,
	}
	relFilepath := "nonexistent/file.txt"
	quartoDirPath := "../utils/test_files/good_repository_setup"

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(50))
	mockFilesystem.EXPECT().CheckoutBranch("10").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return(quartoDirPath)
	mockDiscussionContainerRepository.EXPECT().GetByID(uint(0)).Return(&models.DiscussionContainer{}, nil)
	mockClosedBranchRepository.EXPECT().Query(&models.ClosedBranch{ProjectPostID: 20})

	absFilepath, err := branchService.GetFileFromProject(10, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", absFilepath)
}

func TestCloseBranchButDontMarkProjectPostAsRevisionNeeded(t *testing.T) {
	beforeEachBranch(t)

	postID := uint(5)
	projectPostID := uint(10)
	postDiscussionContainerID := uint(15)

	branchID := uint(15)
	reviewingMemberID := uint(1)

	reviewingMember := &models.Member{
		Model: gorm.Model{ID: reviewingMemberID},
	}

	branch := &models.Branch{
		Model:                     gorm.Model{ID: branchID},
		ProjectPostID:             &projectPostID,
		BranchOverallReviewStatus: models.BranchOpenForReview,
	}

	post := &models.Post{
		Model:                 gorm.Model{ID: postID},
		DiscussionContainer:   models.DiscussionContainer{Model: gorm.Model{ID: postDiscussionContainerID}},
		DiscussionContainerID: postDiscussionContainerID,
	}

	projectPost := &models.ProjectPost{
		Model:  gorm.Model{ID: projectPostID},
		Post:   *post,
		PostID: postID,
		OpenBranches: []*models.Branch{
			branch,
		},
		ClosedBranches:   []*models.ClosedBranch{},
		PostReviewStatus: models.Reviewed,
	}

	reviewCreationForm := forms.ReviewCreationForm{
		BranchID:             branchID,
		BranchReviewDecision: "rejected",
		Feedback:             "ur grammar is bad",
	}

	// Setup mock getters
	mockMemberRepository.EXPECT().GetByID(reviewingMemberID).Return(reviewingMember, nil).Times(1)
	mockBranchRepository.EXPECT().GetByID(branchID).Return(branch, nil).Times(1)
	mockProjectPostRepository.EXPECT().GetByID(projectPostID).Return(projectPost, nil).Times(1)
	mockDiscussionContainerRepository.EXPECT().GetByID(postDiscussionContainerID).Return(&post.DiscussionContainer, nil).Times(1)

	// Setup mock creates & updates
	mockBranchRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
	mockBranchRepository.EXPECT().Update(gomock.Any()).Return(branch, nil).Times(1)
	mockBranchReviewRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	// Setup mock queries
	mockClosedBranchRepository.EXPECT().Query(gomock.Any()).Return([]*models.ClosedBranch{}, nil).Times(1)

	// Use argument capture to extract the project post passed to the repo update method
	var capturedProjectPost *models.ProjectPost

	mockProjectPostRepository.EXPECT().Update(gomock.Any()).Do(func(arg *models.ProjectPost) {
		capturedProjectPost = arg
	})

	// Function under test
	_, err := branchService.CreateReview(reviewCreationForm, &models.Member{Model: gorm.Model{ID: reviewingMemberID}})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, models.Reviewed, capturedProjectPost.PostReviewStatus)
}

func TestCreateReviewFailsWhenAlreadyReviewed(t *testing.T) {
	beforeEachBranch(t)

	// Setup data
	branchID := uint(5)
	branch := &models.Branch{
		BranchOverallReviewStatus: models.BranchRejected,
	}

	// Setup mocks
	mockBranchRepository.EXPECT().GetByID(branchID).Return(branch, nil).Times(1)

	reviewCreationForm := forms.ReviewCreationForm{
		BranchID: branchID,
	}

	// Function under test
	_, err := branchService.CreateReview(reviewCreationForm, &models.Member{})
	assert.NotNil(t, err)
}
