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

	// Setup SUT
	projectPostService = ProjectPostService{
		ProjectPostRepository:   projectPostRepositoryMock,
		MemberRepository:        memberRepositoryMock,
		PostCollaboratorService: postCollaboratorServiceMock,
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
			PostType:            tags.Project,
			ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
		},
		CompletionStatus:   tags.Ongoing,
		FeedbackPreference: tags.FormalFeedback,
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{Member: memberA, CollaborationType: models.Author},
		{Member: memberB, CollaborationType: models.Author},
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
			PostType:            tags.Project,
			ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
			DiscussionContainer: models.DiscussionContainer{
				Discussions: []*models.Discussion{},
			},
		},
		OpenBranches:        []*models.Branch{},
		ClosedBranches:      []*models.ClosedBranch{},
		CompletionStatus:    tags.Ongoing,
		FeedbackPreference:  tags.FormalFeedback,
		PostReviewStatusTag: tags.Open,
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
			PostType:            tags.Project,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   tags.Completed,
		FeedbackPreference: tags.FormalFeedback,
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{}, true, models.Author).Return([]*models.PostCollaborator{}, nil).Times(1)
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
			PostType:            tags.Question,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   tags.Idea,
		FeedbackPreference: tags.Discussion,
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

// When creating a collaborator list fails, project post creation should fail.
func TestCreateProjectPostCollaboratorsFail(t *testing.T) {
	projectPostServiceSetup(t)
	t.Cleanup(projectPostServiceTeardown)

	projectPostCreationForm := forms.ProjectPostCreationForm{
		PostCreationForm: forms.PostCreationForm{
			AuthorMemberIDs:     []uint{10, 15},
			Title:               "",
			Anonymous:           false,
			PostType:            tags.Project,
			ScientificFieldTags: []tags.ScientificField{},
		},
		CompletionStatus:   tags.Idea,
		FeedbackPreference: tags.Discussion,
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
