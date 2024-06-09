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

var postRepository *mocks.MockModelRepositoryInterface[*models.Post]
var memberRepository *mocks.MockModelRepositoryInterface[*models.Member]
var postService PostService

var memberA, memberB, memberC models.Member

func postServiceSetup(t *testing.T) {
	t.Helper()

	// Mock database repositories
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	postRepository = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)
	memberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)

	// Create post service
	postService = PostService{
		PostRepository:   postRepository,
		MemberRepository: memberRepository,
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

	memberRepository.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).AnyTimes()
	memberRepository.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).AnyTimes()
	memberRepository.EXPECT().GetByID(memberC.ID).Return(&memberC, nil).AnyTimes()
	memberRepository.EXPECT().GetByID(uint(0)).Return(nil, fmt.Errorf("member does not exist")).AnyTimes()
}

func postServiceTeardown() {

}

func TestCreatePostGoodWeather(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// The input we will be sending to the function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs: []uint{memberA.ID, memberB.ID},
		Title:           "My Awesome Question",
		Anonymous:       false,
		PostType:        tags.Question,
		ScientificFieldTags: []tags.ScientificField{
			tags.Mathematics, tags.ComputerScience,
		},
	}

	// What we expect the database to receive, called by function under test
	postRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

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
		PostType: tags.Question,
		ScientificFieldTags: []tags.ScientificField{
			tags.Mathematics, tags.ComputerScience,
		},
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}

	if !reflect.DeepEqual(createdPost, expectedPost) {
		t.Fatalf("created post:\n%+v\n did not equal expected post:\n%+v\n", createdPost, expectedPost)
	}
}

// Try to create a Post with a member that exists, and one that doesn't
// This should fail / throw an error
func TestCreatePostNonExistingMembers(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:     []uint{memberA.ID, 0},
		Title:               "My Broken Post",
		Anonymous:           false,
		PostType:            tags.Reflection,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}

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
		PostType:        tags.Question,
		ScientificFieldTags: []tags.ScientificField{
			tags.Mathematics, tags.ComputerScience,
		},
	}

	// What we expect the database to receive, called by function under test
	postRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if err != nil {
		t.Fatalf("creating a post failed: %s", err)
	}

	expectedPost := models.Post{
		Collaborators: []*models.PostCollaborator{},
		Title:         "My Awesome Question",
		PostType:      tags.Question,
		ScientificFieldTags: []tags.ScientificField{
			tags.Mathematics, tags.ComputerScience,
		},
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
	}

	if !reflect.DeepEqual(*createdPost, expectedPost) {
		t.Fatalf("created post:\n%+v\n did not equal expected post:\n%+v\n", *createdPost, expectedPost)
	}
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
		PostType:            tags.Reflection,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}

	postRepository.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("oh no")).Times(1)

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
		PostType:            tags.Project,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}

	postRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if createdPost != nil {
		t.Fatalf("created post:\n%+v\nshould have been nil", createdPost)
	}

	if err == nil {
		t.Fatalf("creating project post using CreatePost should have thrown error")
	}
}
