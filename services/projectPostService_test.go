package services

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var projectPostService ProjectPostService

func projectPostServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	memberA = models.Member{
		Model: gorm.Model{ID: 5},
	}
	memberB = models.Member{
		Model: gorm.Model{ID: 10},
	}

	// Create mocks
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	mockClosedBranchRepository = mocks.NewMockModelRepositoryInterface[*models.ClosedBranch](mockCtrl)
	mockScientificFieldTagContainerReposiotry = mocks.NewMockModelRepositoryInterface[*models.ScientificFieldTagContainer](mockCtrl)
	mockPostRepository = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)

	mockPostCollaboratorService = mocks.NewMockPostCollaboratorService(mockCtrl)
	mockBranchCollaboratorService = mocks.NewMockBranchCollaboratorService(mockCtrl)
	mockBranchService = mocks.NewMockBranchService(mockCtrl)
	mockTagService = mocks.NewMockTagService(mockCtrl)

	// Setup SUT
	projectPostService = ProjectPostService{
		ClosedBranchRepository:                mockClosedBranchRepository,
		PostRepository:                        mockPostRepository,
		ProjectPostRepository:                 mockProjectPostRepository,
		MemberRepository:                      mockMemberRepository,
		ScientificFieldTagContainerRepository: mockScientificFieldTagContainerReposiotry,
		Filesystem:                            mockFilesystem,
		PostCollaboratorService:               mockPostCollaboratorService,
		BranchCollaboratorService:             mockBranchCollaboratorService,
		BranchService:                         mockBranchService,
		TagService:                            mockTagService,
	}
}

func projectPostServiceTeardown() {

}

func TestCreateProjectPostGoodWeather(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	// Input to function under test

	emptyTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{},
	}
	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:       []uint{memberA.ID, memberB.ID},
			Title:                 "My Awesome Project Post",
			Anonymous:             false,
			PostType:              models.Project,
			ScientificFieldTagIDs: []uint{},
		},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.FormalFeedback,
	}
	expectedPost := &models.Post{
		Collaborators: []*models.PostCollaborator{
			{Member: memberA, CollaborationType: models.Author},
			{Member: memberB, CollaborationType: models.Author},
		},
		Title:                       "My Awesome Project Post",
		PostType:                    models.Project,
		ScientificFieldTagContainer: *emptyTagContainer,
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
		RenderStatus: models.Success,
	}

	// Setup mock function return values
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{Member: memberA, CollaborationType: models.Author},
		{Member: memberB, CollaborationType: models.Author},
	}, nil).Times(1)

	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{memberA.ID, memberB.ID}, false).Return([]*models.BranchCollaborator{
		{Member: memberA}, {Member: memberB},
	}, nil).Times(1)

	mockProjectPostRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(0))
	mockFilesystem.EXPECT().CreateRepository().Return(nil)
	mockFilesystem.EXPECT().CreateBranch("0").Return(nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)
	mockScientificFieldTagContainerReposiotry.EXPECT().Create(emptyTagContainer).Return(nil)
	mockPostRepository.EXPECT().Create(expectedPost).Return(nil)

	// Function under test
	createdProjectPost, err404, err500 := projectPostService.CreateProjectPost(&projectPostCreationForm)
	assert.Nil(t, err404)
	assert.Nil(t, err500)

	expectedProjectPost := &models.ProjectPost{
		Post: *expectedPost,
		OpenBranches: []*models.Branch{
			{
				Collaborators: []*models.BranchCollaborator{
					{Member: memberA}, {Member: memberB},
				},
				Reviews: []*models.BranchReview{},
				DiscussionContainer: models.DiscussionContainer{
					Discussions: []*models.Discussion{},
				},
				BranchTitle:               models.InitialPeerReviewBranchName,
				RenderStatus:              models.Success,
				BranchOverallReviewStatus: models.BranchOpenForReview,
			},
		},
		ClosedBranches:            []*models.ClosedBranch{},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.FormalFeedback,
		PostReviewStatus:          models.Open,
	}

	assert.Equal(t, expectedProjectPost, createdProjectPost)
}

// When database creation fails, project post creation should fail.
func TestCreateProjectPostDatabaseFailure(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	emptyTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{},
	}
	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:       []uint{},
			Title:                 "My Broken Project Post",
			Anonymous:             true,
			PostType:              models.Project,
			ScientificFieldTagIDs: []uint{},
		},
		ProjectCompletionStatus:   models.Completed,
		ProjectFeedbackPreference: models.FormalFeedback,
	}
	expectedPost := &models.Post{
		Collaborators:               []*models.PostCollaborator{},
		Title:                       "My Broken Project Post",
		PostType:                    models.Project,
		ScientificFieldTagContainer: *emptyTagContainer,
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
		RenderStatus: models.Success,
	}

	// Setup mock function return values
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{}, true, models.Author).Return([]*models.PostCollaborator{}, nil).Times(1)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{}, true).Return([]*models.BranchCollaborator{}, nil).Times(1)
	mockProjectPostRepository.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("oh no")).Times(1)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)
	mockScientificFieldTagContainerReposiotry.EXPECT().Create(emptyTagContainer).Return(nil)
	mockPostRepository.EXPECT().Create(expectedPost).Return(nil)

	// Function under test
	createdProjectPost, err404, err500 := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	assert.Nil(t, err404)
	assert.NotNil(t, err500)
}

// Creating a project post must use the correct post type, otherwise creation should fail.
func TestCreateProjectPostWrongPostType(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:       []uint{},
			Title:                 "",
			Anonymous:             true,
			PostType:              models.Question,
			ScientificFieldTagIDs: []uint{},
		},
		ProjectCompletionStatus:   models.Idea,
		ProjectFeedbackPreference: models.DiscussionFeedback,
	}

	// Function under test
	createdProjectPost, err404, err500 := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

// When creating a post collaborator list fails, project post creation should fail.
func TestCreateProjectPostCollaboratorsFail(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:       []uint{10, 15},
			Title:                 "",
			Anonymous:             false,
			PostType:              models.Project,
			ScientificFieldTagIDs: []uint{},
		},
		ProjectCompletionStatus:   models.Idea,
		ProjectFeedbackPreference: models.DiscussionFeedback,
	}

	// Setup mock function return values
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{10, 15}, false, models.Author).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	createdProjectPost, err404, err500 := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

// When creating a branch collaborator list fails, project post creation should fail.
func TestCreateProjectBranchCollaboratorsFail(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	emptyTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{},
	}
	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:       []uint{memberA.ID, memberB.ID},
			Title:                 "",
			Anonymous:             false,
			PostType:              models.Project,
			ScientificFieldTagIDs: []uint{},
		},
		ProjectCompletionStatus:   models.Idea,
		ProjectFeedbackPreference: models.DiscussionFeedback,
	}
	expectedPost := &models.Post{
		Collaborators: []*models.PostCollaborator{
			{Member: memberA, CollaborationType: models.Author},
			{Member: memberB, CollaborationType: models.Author},
		},
		Title:                       "",
		PostType:                    models.Project,
		ScientificFieldTagContainer: *emptyTagContainer,
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
		RenderStatus: models.Success,
	}

	// Setup mock function return values
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{Member: memberA, CollaborationType: models.Author},
		{Member: memberB, CollaborationType: models.Author},
	}, nil).Times(1)
	mockBranchCollaboratorService.EXPECT().MembersToBranchCollaborators([]uint{memberA.ID, memberB.ID}, false).Return(nil, fmt.Errorf("oh no")).Times(1)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)
	mockScientificFieldTagContainerReposiotry.EXPECT().Create(emptyTagContainer).Return(nil)
	mockPostRepository.EXPECT().Create(expectedPost).Return(nil)

	// Function under test
	createdProjectPost, err404, err500 := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	assert.NotNil(t, err404)
	assert.Nil(t, err500)
}

func TestGetProjectPost(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	databasePost := &models.ProjectPost{
		Model: gorm.Model{ID: 10},
		Post: models.Post{
			Model:         gorm.Model{ID: 20},
			Collaborators: []*models.PostCollaborator{},
			Title:         "My Awesome Project Post",
			PostType:      models.Project,
			ScientificFieldTagContainer: models.ScientificFieldTagContainer{
				ScientificFieldTags: []*models.ScientificFieldTag{},
			},
			DiscussionContainer: models.DiscussionContainer{
				Model:       gorm.Model{ID: 1},
				Discussions: []*models.Discussion{},
			},
			DiscussionContainerID: 1,
		},
		PostID:                    20,
		OpenBranches:              []*models.Branch{{Model: gorm.Model{ID: 25}}},
		ClosedBranches:            []*models.ClosedBranch{},
		ProjectCompletionStatus:   models.Ongoing,
		ProjectFeedbackPreference: models.FormalFeedback,
		PostReviewStatus:          models.Open,
	}

	mockProjectPostRepository.EXPECT().GetByID(uint(10)).Return(databasePost, nil).Times(1)

	// Function under test
	fetchedPost, err := projectPostService.GetProjectPost(10)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fetchedPost, databasePost) {
		t.Fatalf("fetched post\n%+v\nshould have equaled expected post\n%+v", fetchedPost, databasePost)
	}
}

func TestFilterAllProjectPosts(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	page := 1
	size := 2

	// For this test, we leave the form empty - we want all posts!
	form := forms.ProjectPostFilterForm{}

	// Setup mock function return values
	mockProjectPostRepository.EXPECT().QueryPaginated(page, size, gomock.Any()).Return([]*models.ProjectPost{
		{Model: gorm.Model{ID: 2}},
		{Model: gorm.Model{ID: 3}},
		{Model: gorm.Model{ID: 6}},
		{Model: gorm.Model{ID: 10}},
	}, nil).Times(1)

	// Function under test
	fetchedPostIDs, err := projectPostService.Filter(page, size, form)
	if err != nil {
		t.Fatal(err)
	}

	expectedPostIDs := []uint{2, 3, 6, 10}

	if !reflect.DeepEqual(fetchedPostIDs, expectedPostIDs) {
		t.Fatalf("fetched post IDs\n%+v\nshould have equaled expected post IDs\n%+v", fetchedPostIDs, expectedPostIDs)
	}
}

func TestFilterProjectPostsFailed(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	mockProjectPostRepository.EXPECT().QueryPaginated(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := projectPostService.Filter(1, 10, forms.ProjectPostFilterForm{})

	if err == nil {
		t.Fatal("post filtering should have failed")
	}
}

func TestGetDiscussionContainersFromMergeHistory(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	var projectPostID uint = 5

	databaseProjectPost := &models.ProjectPost{
		Post: models.Post{
			DiscussionContainer: models.DiscussionContainer{
				Model: gorm.Model{ID: 45},
			},
			DiscussionContainerID: 45,
		},
		ClosedBranches: []*models.ClosedBranch{
			{
				Model: gorm.Model{ID: 22},
				Branch: models.Branch{
					Model: gorm.Model{ID: 99},
					DiscussionContainer: models.DiscussionContainer{
						Model: gorm.Model{ID: 54},
					},
					DiscussionContainerID: 54,
				},
				BranchID:             99,
				BranchReviewDecision: models.Approved,
			},
		},
	}

	// Setup mock function return values
	mockProjectPostRepository.EXPECT().GetByID(projectPostID).Return(databaseProjectPost, nil).Times(1)
	mockClosedBranchRepository.EXPECT().Query(gomock.Any()).Return([]*models.ClosedBranch{
		{
			Model: gorm.Model{ID: 22},
			Branch: models.Branch{
				Model: gorm.Model{ID: 99},
				DiscussionContainer: models.DiscussionContainer{
					Model: gorm.Model{ID: 54},
				},
				DiscussionContainerID: 54,
			},
			BranchID:             99,
			BranchReviewDecision: models.Approved,
		},
	}, nil).Times(1)

	// Function under test
	discussionContainerHistory, err := projectPostService.GetDiscussionContainersFromMergeHistory(projectPostID)
	if err != nil {
		t.Fatal(err)
	}

	expectedDiscussionContainerHistory := &models.DiscussionContainerProjectHistoryDTO{
		CurrentDiscussionContainerID: 45,
		MergedBranchDiscussionContainers: []models.DiscussionContainerWithBranchDTO{
			{
				DiscussionContainerID: 54,
				ClosedBranchID:        22,
			},
		},
	}

	if !reflect.DeepEqual(discussionContainerHistory, expectedDiscussionContainerHistory) {
		t.Fatalf("discussion container history\n%+v\ndid not equal expected discussion container history\n%+v", discussionContainerHistory, expectedDiscussionContainerHistory)
	}
}

func TestGetDiscussionContainersFromMergeHistoryPostNotFound(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	mockProjectPostRepository.EXPECT().GetByID(uint(50)).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := projectPostService.GetDiscussionContainersFromMergeHistory(50)

	if err == nil {
		t.Fatal("getting discussion container history should have returned error")
	}
}
