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
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
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
	mockMemberService = mocks.NewMockMemberService(mockCtrl)
	mockRenderService = mocks.NewMockRenderService(mockCtrl)
	mockBranchRepository = mocks.NewMockModelRepositoryInterface[*models.Branch](mockCtrl)
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockReviewRepository = mocks.NewMockModelRepositoryInterface[*models.Review](mockCtrl)
	mockBranchCollaboratorRepository = mocks.NewMockModelRepositoryInterface[*models.BranchCollaborator](mockCtrl)
	mockDiscussionContainerRepository = mocks.NewMockModelRepositoryInterface[*models.DiscussionContainer](mockCtrl)
	mockDiscussionRepository = mocks.NewMockModelRepositoryInterface[*models.Discussion](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)

	// Create branch service
	branchService = BranchService{
		MemberService:                 mockMemberService,
		RenderService:                 mockRenderService,
		BranchRepository:              mockBranchRepository,
		ProjectPostRepository:         mockProjectPostRepository,
		ReviewRepository:              mockReviewRepository,
		BranchCollaboratorRepository:  mockBranchCollaboratorRepository,
		DiscussionContainerRepository: mockDiscussionContainerRepository,
		DiscussionRepository:          mockDiscussionRepository,
		Filesystem:                    mockFilesystem,
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
	projectPost.PostID = 10
	collaborator := &models.BranchCollaborator{MemberID: 12}
	expectedBranch := &models.Branch{
		Collaborators: []*models.BranchCollaborator{collaborator},
		ProjectPostID: 10,
	}
	outputBranch := &models.Branch{
		Model:         gorm.Model{ID: 15},
		Collaborators: []*models.BranchCollaborator{collaborator},
		ProjectPostID: 10,
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockBranchRepository.EXPECT().Create(expectedBranch).DoAndReturn(
		func(expectedBranch *models.Branch) error {
			expectedBranch.ID = 15
			return nil
		})
	mockFilesystem.EXPECT().CheckoutDirectory(uint(10))
	mockFilesystem.EXPECT().CreateBranch("15").Return(nil)

	branch, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		Collaborators: []*models.BranchCollaborator{collaborator},
		ProjectPostID: 10,
	})

	assert.Nil(t, err404)
	assert.Nil(t, err500)
	assert.Equal(t, outputBranch, &branch)
}

func TestCreateBranchNoProjectPost(t *testing.T) {
	beforeEachBranch(t)

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, errors.New("failed"))

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		Collaborators: []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID: 10,
	})

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

func TestCreateBranchFailedToCreate(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	expectedBranch := &models.Branch{
		Collaborators: []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID: 10,
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockBranchRepository.EXPECT().Create(expectedBranch).Return(errors.New("failed"))

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		Collaborators: []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID: 10,
	})

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestCreateBranchFailedGit(t *testing.T) {
	beforeEachBranch(t)

	projectPost.ID = 10
	projectPost.PostID = 10
	expectedBranch := &models.Branch{
		Collaborators: []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID: 10,
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockBranchRepository.EXPECT().Create(expectedBranch).DoAndReturn(
		func(expectedBranch *models.Branch) error {
			expectedBranch.ID = 15
			return nil
		})
	mockFilesystem.EXPECT().CheckoutDirectory(uint(10))
	mockFilesystem.EXPECT().CreateBranch("15").Return(errors.New("failed"))

	_, err404, err500 := branchService.CreateBranch(&forms.BranchCreationForm{
		Collaborators: []*models.BranchCollaborator{{MemberID: 12, BranchID: 11}},
		ProjectPostID: 10,
	})

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

func TestUpdateBranchSuccess(t *testing.T) {
	beforeEachBranch(t)

	input := models.BranchDTO{
		ID:              1,
		NewPostTitle:    "test",
		CollaboratorIDs: []uint{5},
		DiscussionIDs:   []uint{6},
		ProjectPostID:   10,
	}
	collaborator := &models.BranchCollaborator{Model: gorm.Model{ID: 20}}
	discussion := &models.Discussion{Model: gorm.Model{ID: 21}}
	expected := &models.Branch{
		Model:               gorm.Model{ID: 1},
		NewPostTitle:        "test",
		ProjectPostID:       10,
		Collaborators:       []*models.BranchCollaborator{collaborator},
		DiscussionContainer: models.DiscussionContainer{Discussions: []*models.Discussion{discussion}},
	}

	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(5)).Return(collaborator, nil)
	mockDiscussionRepository.EXPECT().GetByID(uint(6)).Return(discussion, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockBranchRepository.EXPECT().Update(expected).Return(expected, nil)

	output, err := branchService.UpdateBranch(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, &output)
}

func TestUpdateBranchNoSuchCollaborator(t *testing.T) {
	beforeEachBranch(t)

	input := models.BranchDTO{
		ID:              1,
		NewPostTitle:    "test",
		CollaboratorIDs: []uint{5},
	}

	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(5)).Return(&models.BranchCollaborator{MemberID: 19}, errors.New("failed"))

	_, err := branchService.UpdateBranch(input)
	assert.NotNil(t, err)
}

func TestUpdateNoSuchDiscussion(t *testing.T) {
	beforeEachBranch(t)

	input := models.BranchDTO{
		ID:              1,
		NewPostTitle:    "test",
		CollaboratorIDs: []uint{5},
		DiscussionIDs:   []uint{6},
	}

	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(5)).Return(&models.BranchCollaborator{Model: gorm.Model{ID: 20}}, nil)
	mockDiscussionRepository.EXPECT().GetByID(uint(6)).Return(&models.Discussion{Model: gorm.Model{ID: 21}}, errors.New("failed"))

	_, err := branchService.UpdateBranch(input)
	assert.NotNil(t, err)
}

func TestUpdateBranchNoSuchProjectPost(t *testing.T) {
	beforeEachBranch(t)

	input := models.BranchDTO{
		ID:              1,
		NewPostTitle:    "test",
		CollaboratorIDs: []uint{5},
		DiscussionIDs:   []uint{6},
		ProjectPostID:   10,
	}

	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(5)).Return(&models.BranchCollaborator{Model: gorm.Model{ID: 20}}, nil)
	mockDiscussionRepository.EXPECT().GetByID(uint(6)).Return(&models.Discussion{Model: gorm.Model{ID: 21}}, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, errors.New("failed"))

	_, err := branchService.UpdateBranch(input)
	assert.NotNil(t, err)
}

func TestUpdateBranchFailedUpdate(t *testing.T) {
	beforeEachBranch(t)

	input := models.BranchDTO{
		ID:              1,
		NewPostTitle:    "test",
		CollaboratorIDs: []uint{5},
		DiscussionIDs:   []uint{6},
		ProjectPostID:   10,
	}
	collaborator := &models.BranchCollaborator{Model: gorm.Model{ID: 20}}
	discussion := &models.Discussion{Model: gorm.Model{ID: 21}}
	expected := &models.Branch{
		Model:               gorm.Model{ID: 1},
		NewPostTitle:        "test",
		ProjectPostID:       10,
		Collaborators:       []*models.BranchCollaborator{collaborator},
		DiscussionContainer: models.DiscussionContainer{Discussions: []*models.Discussion{discussion}},
	}

	mockBranchCollaboratorRepository.EXPECT().GetByID(uint(5)).Return(collaborator, nil)
	mockDiscussionRepository.EXPECT().GetByID(uint(6)).Return(discussion, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(projectPost, nil)
	mockBranchRepository.EXPECT().Update(expected).Return(expected, errors.New("failed"))

	_, err := branchService.UpdateBranch(input)
	assert.NotNil(t, err)
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
		Model:        gorm.Model{ID: 10},
		NewPostTitle: "title",
		Reviews:      []*models.Review{{BranchDecision: models.Approved}, {BranchDecision: models.Rejected}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)

	decisions, err := branchService.GetReviewStatus(uint(10))
	assert.Nil(t, err)
	assert.Equal(t, []models.BranchDecision{models.Approved, models.Rejected}, decisions)
}

func TestGetReviewStatusFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model:        gorm.Model{ID: 10},
		NewPostTitle: "title",
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
		BranchID:       10,
		MemberID:       11,
		BranchDecision: models.Approved,
	}
	expected := &models.Review{
		// Model:          gorm.Model{ID: 1},
		BranchID:       10,
		Member:         models.Member{Model: gorm.Model{ID: 11}},
		BranchDecision: models.Approved,
	}
	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	newBranch := &models.Branch{
		Model:   gorm.Model{ID: 10},
		Reviews: []*models.Review{expected},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, nil)
	mockReviewRepository.EXPECT().Create(expected).Return(nil)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, nil)

	review, err := branchService.CreateReview(form)
	assert.Nil(t, err)
	assert.Equal(t, expected, &review)
}

func TestCreateReviewFailedGetBranch(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	form := forms.ReviewCreationForm{
		BranchID:       10,
		MemberID:       11,
		BranchDecision: models.Approved,
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
		BranchID:       10,
		MemberID:       11,
		BranchDecision: models.Approved,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, errors.New("failed"))

	_, err := branchService.CreateReview(form)
	assert.NotNil(t, err)
}

func TestCreateReviewFailedCreateReview(t *testing.T) {
	beforeEachBranch(t)

	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:       10,
		MemberID:       11,
		BranchDecision: models.Approved,
	}
	expected := &models.Review{
		// Model:          gorm.Model{ID: 1},
		BranchID:       10,
		Member:         models.Member{Model: gorm.Model{ID: 11}},
		BranchDecision: models.Approved,
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, nil)
	mockReviewRepository.EXPECT().Create(expected).Return(errors.New("failed"))

	_, err := branchService.CreateReview(form)
	assert.NotNil(t, err)
}

func TestCreateReviewFailedUpdateBranch(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model: gorm.Model{ID: 11},
	}
	form := forms.ReviewCreationForm{
		BranchID:       10,
		MemberID:       11,
		BranchDecision: models.Approved,
	}
	expected := &models.Review{
		// Model:          gorm.Model{ID: 1},
		BranchID:       10,
		Member:         models.Member{Model: gorm.Model{ID: 11}},
		BranchDecision: models.Approved,
	}
	branch := &models.Branch{
		Model: gorm.Model{ID: 10},
	}
	newBranch := &models.Branch{
		Model:   gorm.Model{ID: 10},
		Reviews: []*models.Review{expected},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, nil)
	mockReviewRepository.EXPECT().Create(expected).Return(nil)
	mockBranchRepository.EXPECT().Update(newBranch).Return(newBranch, errors.New("failed"))

	_, err := branchService.CreateReview(form)
	assert.NotNil(t, err)
}

func TestMemberCanReviewSuccessTrue(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model:               gorm.Model{ID: 11},
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 20,
	}
	projectPost := &models.ProjectPost{
		Model: gorm.Model{ID: 20},
		Post:  models.Post{ScientificFieldTags: []tags.ScientificField{tags.Mathematics, tags.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, nil)

	canReview, err := branchService.MemberCanReview(10, 11)
	assert.Nil(t, err)
	assert.True(t, canReview)
}

func TestMemberCanReviewSuccessFalse(t *testing.T) {
	beforeEachBranch(t)

	member := &models.Member{
		Model:               gorm.Model{ID: 11},
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}
	branch := &models.Branch{
		Model:         gorm.Model{ID: 10},
		ProjectPostID: 20,
	}
	projectPost := &models.ProjectPost{
		Model: gorm.Model{ID: 20},
		Post:  models.Post{ScientificFieldTags: []tags.ScientificField{tags.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, nil)

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
		Post:  models.Post{ScientificFieldTags: []tags.ScientificField{tags.Mathematics, tags.ComputerScience}},
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
		Post:  models.Post{ScientificFieldTags: []tags.ScientificField{tags.Mathematics, tags.ComputerScience}},
	}

	mockBranchRepository.EXPECT().GetByID(uint(10)).Return(branch, nil)
	mockProjectPostRepository.EXPECT().GetByID(uint(20)).Return(projectPost, nil)
	mockMemberService.EXPECT().GetMember(uint(11)).Return(member, errors.New("failed"))

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
	mockRenderService.EXPECT().Render(branch)

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
