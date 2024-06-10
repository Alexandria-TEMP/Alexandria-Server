package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var postService PostService

func postServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	postRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)
	memberRepositoryMock = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	postCollaboratorServiceMock = mocks.NewMockPostCollaboratorService(mockCtrl)

	// Setup SUT
	postService = PostService{
		PostRepository:          postRepositoryMock,
		MemberRepository:        memberRepositoryMock,
		PostCollaboratorService: postCollaboratorServiceMock,
	}

	// Setup members in the repository
	memberA = models.Member{
		Model: gorm.Model{ID: 5},
	}

	memberB = models.Member{
		Model: gorm.Model{ID: 10},
	}

	memberC = models.Member{
		Model: gorm.Model{ID: 12},
	}

	memberRepositoryMock.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).AnyTimes()
	memberRepositoryMock.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).AnyTimes()
	memberRepositoryMock.EXPECT().GetByID(memberC.ID).Return(&memberC, nil).AnyTimes()
	memberRepositoryMock.EXPECT().GetByID(uint(0)).Return(nil, fmt.Errorf("member does not exist")).AnyTimes()
}

func postServiceTeardown() {

}

func TestCreatePostGoodWeather(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// The input we will be sending to the function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:     []uint{memberA.ID, memberB.ID},
		Title:               "My Awesome Question",
		Anonymous:           false,
		PostType:            models.Question,
		ScientificFieldTags: []*tags.ScientificFieldTag{},
	}

	// Setup mock function return values
	postRepositoryMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{
			Member:            memberA,
			CollaborationType: models.Author,
		},
		{
			Member:            memberB,
			CollaborationType: models.Author,
		},
	}, nil).Times(1)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if err != nil {
		t.Fatalf("creating a post failed: %s", err)
	}

	expectedPost := &models.Post{
		Collaborators: []*models.PostCollaborator{
			{
				Member:            memberA,
				CollaborationType: models.Author,
			},
			{
				Member:            memberB,
				CollaborationType: models.Author,
			},
		},
		Title:    "My Awesome Question",
		PostType: models.Question,
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}
	assert.Equal(t, createdPost, expectedPost)
}

// Try to create a Post where the PostCollaboratorService returns an error. Should fail.
func TestCreatePostNonExistingMembers(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:     []uint{memberA.ID, memberB.ID},
		Title:               "My Broken Post",
		Anonymous:           false,
		PostType:            models.Reflection,
		ScientificFieldTags: ScientificFieldTags: []*tags.ScientificFieldTag{},
	}

	// Setup mock function return values
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if createdPost != nil {
		t.Fatalf("created post:\n%+v\nshould have been nil", createdPost)
	}

	if err == nil {
		t.Fatalf("creating post with invalid member should have thrown error")
	}
}

// Creating a post with anonymity should give an empty list of collaborators,
// even if author member IDs are given!
func TestCreatePostWithAnonymity(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// The input we will be sending to the function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs: []uint{memberA.ID, memberB.ID},
		Title:           "My Awesome Question",
		Anonymous:       true,
		PostType:        models.Question,
		ScientificFieldTags: []*tags.ScientificFieldTag{},
	}

	// Setup mock function return values
	postRepositoryMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, true, models.Author).Return([]*models.PostCollaborator{}, nil)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if err != nil {
		t.Fatalf("creating a post failed: %s", err)
	}

	expectedPost := models.Post{
		Collaborators: []*models.PostCollaborator{},
		Title:         "My Awesome Question",
		PostType:      models.Question,
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}
	assert.Equal(t, *createdPost, expectedPost)
}

// If the database creation fails, creating a post should fail
func TestCreatePostDatabaseFailure(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:     []uint{memberA.ID, memberC.ID},
		Title:               "My Post That Shall Fail",
		Anonymous:           false,
		PostType:            models.Reflection,
		ScientificFieldTags: []*tags.ScientificFieldTag{},
	}

	postRepositoryMock.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("oh no")).Times(1)
	postCollaboratorServiceMock.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberC.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{
			Member:            memberA,
			CollaborationType: models.Author,
		},
		{
			Member:            memberC,
			CollaborationType: models.Author,
		},
	}, nil)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if createdPost != nil {
		t.Fatalf("created post:\n%+v\nshould have been nil", createdPost)
	}

	if err == nil {
		t.Fatalf("creating post causing database failure should have thrown error")
	}
}

// Creating a ProjectPost should not work with the CreatePost method,
// because it requires extra data.
func TestCreatePostWithBadPostType(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:     []uint{memberA.ID, memberB.ID, memberC.ID},
		Title:               "My Faulty Project Post",
		Anonymous:           false,
		PostType:            models.Project,
		ScientificFieldTags: []*tags.ScientificFieldTag{},
	}

	postRepositoryMock.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if createdPost != nil {
		t.Fatalf("created post:\n%+v\nshould have been nil", createdPost)
	}

	if err == nil {
		t.Fatalf("creating project post using CreatePost should have thrown error")
	}
}

func TestGetPost(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	databasePost := &models.Post{
		Model:               gorm.Model{ID: 5},
		Collaborators:       []*models.PostCollaborator{},
		Title:               "Hello, world!",
		PostType:            models.Project,
		ScientificFieldTags: []tags.ScientificField{},
		DiscussionContainer: models.DiscussionContainer{
			Model:       gorm.Model{ID: 6},
			Discussions: []*models.Discussion{},
		},
		DiscussionContainerID: 6,
	}

	postRepositoryMock.EXPECT().GetByID(uint(10)).Return(databasePost, nil).Times(1)

	// Function under test
	fetchedPost, err := postService.GetPost(10)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fetchedPost, databasePost) {
		t.Fatalf("fetched post\n%+v\nshould have equaled expected post\n%+v", fetchedPost, databasePost)
	}
}

func TestFilterAllPosts(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	page := 1
	size := 2

	// For this test, we leave the form empty - we want all posts!
	form := forms.FilterForm{}

	// Setup mock function return values
	postRepositoryMock.EXPECT().QueryPaginated(page, size, gomock.Any()).Return([]*models.Post{
		{Model: gorm.Model{ID: 2}},
		{Model: gorm.Model{ID: 3}},
		{Model: gorm.Model{ID: 6}},
		{Model: gorm.Model{ID: 10}},
	}, nil).Times(1)

	// Function under test
	fetchedPostIDs, err := postService.Filter(page, size, form)
	if err != nil {
		t.Fatal(err)
	}

	expectedPostIDs := []uint{2, 3, 6, 10}

	if !reflect.DeepEqual(fetchedPostIDs, expectedPostIDs) {
		t.Fatalf("fetched post IDs\n%+v\nshould have equaled expected post IDs\n%+v", fetchedPostIDs, expectedPostIDs)
	}
}

func TestFilterFailed(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postRepositoryMock.EXPECT().QueryPaginated(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := postService.Filter(1, 10, forms.FilterForm{})

	if err == nil {
		t.Fatal("post filtering should have failed")
	}
}
