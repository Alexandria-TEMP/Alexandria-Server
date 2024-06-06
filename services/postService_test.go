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
		AuthorMemberIDs: []uint{memberA.ID, memberB.ID},
		Title:           "My Awesome Question",
		Anonymous:       false,
		PostType:        tags.Question,
		ScientificFieldTags: []tags.ScientificField{
			tags.Mathematics, tags.ComputerScience,
		},
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

// Try to create a Post where the PostCollaboratorService returns an error. Should fail.
func TestCreatePostNonExistingMembers(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:     []uint{memberA.ID, memberB.ID},
		Title:               "My Broken Post",
		Anonymous:           false,
		PostType:            tags.Reflection,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
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
		PostType:        tags.Question,
		ScientificFieldTags: []tags.ScientificField{
			tags.Mathematics, tags.ComputerScience,
		},
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
		PostType:            tags.Project,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
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
