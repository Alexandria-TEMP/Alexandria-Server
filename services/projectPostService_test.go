package services

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
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
	projectPostRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	memberRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)

	postCollaboratorServiceMock = mocks.NewMockPostCollaboratorService(mockCtrl)
	branchCollaboratorServiceMock = mocks.NewMockBranchCollaboratorService(mockCtrl)

	// Setup SUT
	projectPostService = ProjectPostService{
		ProjectPostRepository:     projectPostRepositoryMock,
		MemberRepository:          memberRepositoryMock,
		PostCollaboratorService:   postCollaboratorServiceMock,
		BranchCollaboratorService: branchCollaboratorServiceMock,
	}
}

func projectPostServiceTeardown() {

}

func TestCreateProjectPostGoodWeather(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	// Input to function under test
	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:     []uint{memberA.ID, memberB.ID},
			Title:               "My Awesome Project Post",
			Anonymous:           false,
			PostType:            models.Project,
			ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
		},
		CompletionStatus:   models.Ongoing,
		FeedbackPreference: models.FormalFeedback,
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{Member: memberA, CollaborationType: models.Author},
		{Member: memberB, CollaborationType: models.Author},
	}, nil).Times(1)

	branchCollaboratorServiceMock.EXPECT().MembersToBranchCollaborators([]uint{memberA.ID, memberB.ID}, false).Return([]*models.BranchCollaborator{
		{Member: memberA}, {Member: memberB},
	}, nil).Times(1)

	projectPostRepositoryMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	// Function under test
	createdProjectPost, err := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if err != nil {
		t.Fatalf("creating project post failed, reason: %s", err)
	}

	expectedProjectPost := &models.ProjectPost{
		Post: models.Post{
			Collaborators: []*models.PostCollaborator{
				{Member: memberA, CollaborationType: models.Author},
				{Member: memberB, CollaborationType: models.Author},
			},
			Title:               "My Awesome Project Post",
			PostType:            models.Project,
			ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
			DiscussionContainer: models.DiscussionContainer{
				Discussions: []*models.Discussion{},
			},
		},
		OpenBranches: []*models.Branch{
			{
				NewPostTitle:            "My Awesome Project Post",
				UpdatedCompletionStatus: models.Ongoing,
				UpdatedScientificFields: []tags.ScientificField{tags.Mathematics},
				Collaborators: []*models.BranchCollaborator{
					{Member: memberA}, {Member: memberB},
				},
				Reviews: []*models.BranchReview{},
				DiscussionContainer: models.DiscussionContainer{
					Discussions: []*models.Discussion{},
				},
				BranchTitle:        models.InitialPeerReviewBranchName,
				RenderStatus:       models.Pending,
				BranchReviewStatus: models.BranchOpenForReview,
			},
		},
		ClosedBranches:     []*models.ClosedBranch{},
		CompletionStatus:   models.Ongoing,
		FeedbackPreference: models.FormalFeedback,
		PostReviewStatus:   models.Open,
	}

	if !reflect.DeepEqual(createdProjectPost, expectedProjectPost) {
		t.Fatalf("created project post:\n%+v\ndid not equal expected project post:\n%+v\n",
			createdProjectPost, expectedProjectPost)
	}
}

// When database creation fails, project post creation should fail.
func TestCreateProjectPostDatabaseFailure(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:     []uint{},
			Title:               "My Broken Project Post",
			Anonymous:           true,
			PostType:            models.Project,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   models.Completed,
		FeedbackPreference: models.FormalFeedback,
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{}, true, models.Author).Return([]*models.PostCollaborator{}, nil).Times(1)
	branchCollaboratorServiceMock.EXPECT().MembersToBranchCollaborators([]uint{}, true).Return([]*models.BranchCollaborator{}, nil).Times(1)
	projectPostRepositoryMock.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("oh no")).Times(1)

	// Function under test
	createdProjectPost, err := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	if err == nil {
		t.Fatal("project post creation should have thrown error")
	}
}

// Creating a project post must use the correct post type, otherwise creation should fail.
func TestCreateProjectPostWrongPostType(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:     []uint{},
			Title:               "",
			Anonymous:           true,
			PostType:            models.Question,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   models.Idea,
		FeedbackPreference: models.DiscussionFeedback,
	}

	// Function under test
	createdProjectPost, err := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	if err == nil {
		t.Fatal("project post creation should have thrown error")
	}
}

// When creating a post collaborator list fails, project post creation should fail.
func TestCreateProjectPostCollaboratorsFail(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:     []uint{10, 15},
			Title:               "",
			Anonymous:           false,
			PostType:            models.Project,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   models.Idea,
		FeedbackPreference: models.DiscussionFeedback,
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{10, 15}, false, models.Author).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	createdProjectPost, err := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	if err == nil {
		t.Fatal("project post creation should have thrown error")
	}
}

// When creating a branch collaborator list fails, project post creation should fail.
func TestCreateProjectBranchCollaboratorsFail(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:     []uint{memberA.ID, memberB.ID},
			Title:               "",
			Anonymous:           false,
			PostType:            models.Project,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   models.Idea,
		FeedbackPreference: models.DiscussionFeedback,
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{Member: memberA, CollaborationType: models.Author},
		{Member: memberB, CollaborationType: models.Author},
	}, nil).Times(1)
	branchCollaboratorServiceMock.EXPECT().MembersToBranchCollaborators([]uint{memberA.ID, memberB.ID}, false).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	createdProjectPost, err := projectPostService.CreateProjectPost(&projectPostCreationForm)

	if createdProjectPost != nil {
		t.Fatalf("project post should not have been created:\n%+v", createdProjectPost)
	}

	if err == nil {
		t.Fatal("project post creation should have thrown error")
	}
}

func TestGetProjectPost(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	databasePost := &models.ProjectPost{
		Model: gorm.Model{ID: 10},
		Post: models.Post{
			Model:               gorm.Model{ID: 20},
			Collaborators:       []*models.PostCollaborator{},
			Title:               "My Awesome Project Post",
			PostType:            models.Project,
			ScientificFieldTags: []tags.ScientificField{},
			DiscussionContainer: models.DiscussionContainer{
				Model:       gorm.Model{ID: 1},
				Discussions: []*models.Discussion{},
			},
			DiscussionContainerID: 1,
		},
		PostID:             20,
		OpenBranches:       []*models.Branch{{Model: gorm.Model{ID: 25}}},
		ClosedBranches:     []*models.ClosedBranch{},
		CompletionStatus:   models.Ongoing,
		FeedbackPreference: models.FormalFeedback,
		PostReviewStatus:   models.Open,
	}

	projectPostRepositoryMock.EXPECT().GetByID(uint(10)).Return(databasePost, nil).Times(1)

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
	form := forms.FilterForm{}

	// Setup mock function return values
	projectPostRepositoryMock.EXPECT().QueryPaginated(page, size, gomock.Any()).Return([]*models.ProjectPost{
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

	projectPostRepositoryMock.EXPECT().QueryPaginated(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := projectPostService.Filter(1, 10, forms.FilterForm{})

	if err == nil {
		t.Fatal("post filtering should have failed")
	}
}
