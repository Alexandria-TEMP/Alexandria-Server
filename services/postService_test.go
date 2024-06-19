package services

import (
	"errors"
	"fmt"
	"mime/multipart"
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
var postService PostService

func postServiceSetup(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockPostRepository = mocks.NewMockModelRepositoryInterface[*models.Post](mockCtrl)
	mockProjectPostRepository = mocks.NewMockModelRepositoryInterface[*models.ProjectPost](mockCtrl)
	mockMemberRepository = mocks.NewMockModelRepositoryInterface[*models.Member](mockCtrl)
	mockScientificFieldTagContainerReposiotry = mocks.NewMockModelRepositoryInterface[*models.ScientificFieldTagContainer](mockCtrl)
	mockFilesystem = mocks.NewMockFilesystem(mockCtrl)
	mockPostCollaboratorService = mocks.NewMockPostCollaboratorService(mockCtrl)
	mockRenderService = mocks.NewMockRenderService(mockCtrl)
	mockTagService = mocks.NewMockTagService(mockCtrl)

	// Setup SUT
	postService = PostService{
		PostRepository:                        mockPostRepository,
		ProjectPostRepository:                 mockProjectPostRepository,
		MemberRepository:                      mockMemberRepository,
		ScientificFieldTagContainerRepository: mockScientificFieldTagContainerReposiotry,
		Filesystem:                            mockFilesystem,
		PostCollaboratorService:               mockPostCollaboratorService,
		RenderService:                         mockRenderService,
		TagService:                            mockTagService,
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

	mockMemberRepository.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(memberC.ID).Return(&memberC, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(uint(0)).Return(nil, fmt.Errorf("member does not exist")).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(memberA.ID).Return(&memberA, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(memberB.ID).Return(&memberB, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(memberC.ID).Return(&memberC, nil).AnyTimes()
	mockMemberRepository.EXPECT().GetByID(uint(0)).Return(nil, fmt.Errorf("member does not exist")).AnyTimes()
}

func postServiceTeardown() {

}

func TestCreatePostGoodWeather(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// The input we will be sending to the function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:       []uint{memberA.ID, memberB.ID},
		Title:                 "My Awesome Question",
		Anonymous:             false,
		PostType:              models.Question,
		ScientificFieldTagIDs: []uint{},
	}
	emptyTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{},
	}

	// Setup mock function return values
	mockPostRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{
			Member:            memberA,
			CollaborationType: models.Author,
		},
		{
			Member:            memberB,
			CollaborationType: models.Author,
		},
	}, nil).Times(1)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(0))
	mockFilesystem.EXPECT().CreateRepository().Return(nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)
	mockScientificFieldTagContainerReposiotry.EXPECT().Create(emptyTagContainer).Return(nil)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	assert.Nil(t, err)

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
		Title:                       "My Awesome Question",
		PostType:                    models.Question,
		ScientificFieldTagContainer: *emptyTagContainer,
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
		RenderStatus: models.Success,
	}
	assert.Equal(t, createdPost, expectedPost)
}

// Try to create a Post where the PostCollaboratorService returns an error. Should fail.
func TestCreatePostNonExistingMembers(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:       []uint{memberA.ID, memberB.ID},
		Title:                 "My Broken Post",
		Anonymous:             false,
		PostType:              models.Reflection,
		ScientificFieldTagIDs: []uint{},
	}

	// Setup mock function return values
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return(nil, fmt.Errorf("oh no")).Times(1)
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, false, models.Author).Return(nil, fmt.Errorf("oh no")).Times(1)

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
		AuthorMemberIDs:       []uint{memberA.ID, memberB.ID},
		Title:                 "My Awesome Question",
		Anonymous:             true,
		PostType:              models.Question,
		ScientificFieldTagIDs: []uint{},
	}
	emptyTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{},
	}
	// Setup mock function return values
	mockPostRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberB.ID}, true, models.Author).Return([]*models.PostCollaborator{}, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(0))
	mockFilesystem.EXPECT().CreateRepository().Return(nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)
	mockScientificFieldTagContainerReposiotry.EXPECT().Create(emptyTagContainer).Return(nil)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	assert.Nil(t, err)

	expectedPost := models.Post{
		Collaborators:               []*models.PostCollaborator{},
		Title:                       "My Awesome Question",
		PostType:                    models.Question,
		ScientificFieldTagContainer: *emptyTagContainer,
		DiscussionContainer: models.DiscussionContainer{
			Discussions: []*models.Discussion{},
		},
		RenderStatus: models.Success,
	}
	assert.Equal(t, *createdPost, expectedPost)
}

// If the database creation fails, creating a post should fail
func TestCreatePostDatabaseFailure(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	// Input to function under test
	postCreationForm := forms.PostCreationForm{
		AuthorMemberIDs:       []uint{memberA.ID, memberC.ID},
		Title:                 "My Post That Shall Fail",
		Anonymous:             false,
		PostType:              models.Reflection,
		ScientificFieldTagIDs: []uint{},
	}
	emptyTagContainer := &models.ScientificFieldTagContainer{
		ScientificFieldTags: []*models.ScientificFieldTag{},
	}

	mockPostRepository.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("oh no")).Times(1)
	mockPostCollaboratorService.EXPECT().MembersToPostCollaborators([]uint{memberA.ID, memberC.ID}, false, models.Author).Return([]*models.PostCollaborator{
		{
			Member:            memberA,
			CollaborationType: models.Author,
		},
		{
			Member:            memberC,
			CollaborationType: models.Author,
		},
	}, nil)
	mockTagService.EXPECT().GetTagsFromIDs([]uint{}).Return([]*models.ScientificFieldTag{}, nil)
	mockScientificFieldTagContainerReposiotry.EXPECT().Create(emptyTagContainer).Return(nil)

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
		AuthorMemberIDs:       []uint{memberA.ID, memberB.ID, memberC.ID},
		Title:                 "My Faulty Project Post",
		Anonymous:             false,
		PostType:              models.Project,
		ScientificFieldTagIDs: []uint{},
	}

	mockPostRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
	mockPostRepository.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	// Function under test
	createdPost, err := postService.CreatePost(&postCreationForm)

	if createdPost != nil {
		t.Fatalf("created post:\n%+v\nshould have been nil", createdPost)
	}

	if err == nil {
		t.Fatalf("creating project post using CreatePost should have thrown error")
	}
}

func TestUploadPostSuccess(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	post := &models.Post{
		Model: gorm.Model{ID: 10},
	}
	pendingPost := &models.Post{
		Model:        gorm.Model{ID: 10},
		RenderStatus: models.Pending,
	}
	file := &multipart.FileHeader{
		Filename: "test.zip",
	}

	mockPostRepository.EXPECT().GetByID(uint(10)).Return(post, lock)
	mockFilesystem.EXPECT().LockDirectory(post.ID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(10))
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(nil)
	mockFilesystem.EXPECT().SaveZipFile(gomock.Any(), file).Return(nil)
	mockFilesystem.EXPECT().CreateCommit()
	mockPostRepository.EXPECT().Update(pendingPost).Return(pendingPost, nil)
	mockRenderService.EXPECT().RenderPost(post, lock)

	err := postService.UploadPost(nil, file, 10)
	assert.Nil(t, err)
}

func TestUploadPostFailedGetPost(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	file := &multipart.FileHeader{
		Filename: "test.zip",
	}

	mockPostRepository.EXPECT().GetByID(uint(10)).Return(nil, errors.New("failed"))

	err := postService.UploadPost(nil, file, 10)
	assert.NotNil(t, err)
}

func TestUploadPostFailedCheckoutBranch(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	post := &models.Post{
		Model: gorm.Model{ID: 10},
	}
	file := &multipart.FileHeader{
		Filename: "test.zip",
	}

	mockPostRepository.EXPECT().GetByID(uint(10)).Return(post, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(10))
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(errors.New("failed"))

	err := postService.UploadPost(nil, file, 10)
	assert.NotNil(t, err)
}

func TestUploadPostFailedCleanDir(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	post := &models.Post{
		Model: gorm.Model{ID: 10},
	}
	file := &multipart.FileHeader{
		Filename: "test.zip",
	}

	mockPostRepository.EXPECT().GetByID(uint(10)).Return(post, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(10))
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(errors.New("failed"))

	err := postService.UploadPost(nil, file, 10)
	assert.NotNil(t, err)
}

func TestUploadPostFailedSaveZipFile(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	post := &models.Post{
		Model:        gorm.Model{ID: 10},
		RenderStatus: models.Success,
	}
	file := &multipart.FileHeader{
		Filename: "test.zip",
	}

	mockPostRepository.EXPECT().GetByID(uint(10)).Return(post, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(uint(10))
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().CleanDir().Return(nil)
	mockFilesystem.EXPECT().SaveZipFile(gomock.Any(), file).Return(errors.New("failed"))
	mockPostRepository.EXPECT().Update(gomock.Any()).Return(post, nil)
	mockFilesystem.EXPECT().Reset()

	err := postService.UploadPost(nil, file, 10)
	assert.NotNil(t, err)
	assert.Equal(t, models.Failure, post.RenderStatus)
}

func TestGetMainProjectSuccess(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().LockDirectory(postID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().GetCurrentZipFilePath().Return("../utils/test_files/good_repository_setup/quarto_project.zip")

	filePath, outputLock, err := postService.GetMainProject(postID)
	assert.Nil(t, err)
	assert.Equal(t, "../utils/test_files/good_repository_setup/quarto_project.zip", filePath)
	assert.Equal(t, lock, outputLock)
}

func TestGetMainProjectFailedGetPost(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	mockPostRepository.EXPECT().GetByID(postID).Return(nil, errors.New("failed"))

	filePath, lock, err := postService.GetMainProject(postID)
	assert.NotNil(t, err)
	assert.Nil(t, lock)
	assert.Equal(t, "", filePath)
}

func TestGetMainProjectFailedCheckoutBranch(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().LockDirectory(postID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(errors.New("failed"))

	filePath, outputLock, err := postService.GetMainProject(postID)
	assert.NotNil(t, err)
	assert.Equal(t, "", filePath)
	assert.Equal(t, lock, outputLock)
}

func TestGetMainFiletreeSuccess(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	expectedFileTree := map[string]int64{"file1.txt": 1234, "file2.txt": 5678}
	postID := uint(10)
	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().GetFileTree().Return(expectedFileTree, nil)

	fileTree, err, err2 := postService.GetMainFiletree(postID)
	assert.Nil(t, err)
	assert.Nil(t, err2)
	assert.Equal(t, expectedFileTree, fileTree)
}

func TestGetMainFiletreeFailedGetPost(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	mockPostRepository.EXPECT().GetByID(postID).Return(nil, errors.New("failed"))

	fileTree, err, err2 := postService.GetMainFiletree(postID)
	assert.NotNil(t, err)
	assert.Nil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetMainFiletreeFailedCheckoutBranch(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(errors.New("failed"))

	fileTree, err, err2 := postService.GetMainFiletree(postID)
	assert.NotNil(t, err)
	assert.Nil(t, err2)
	assert.Nil(t, fileTree)
}

func TestGetMainFileFromProjectSuccess(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	relFilepath := "child_dir/test.txt"
	absFilepath := "../utils/test_files/file_tree/child_dir/test.txt"

	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().LockDirectory(postID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return("../utils/test_files/file_tree")

	resultFilepath, outputLock, err := postService.GetMainFileFromProject(postID, relFilepath)
	assert.Nil(t, err)
	assert.Equal(t, absFilepath, resultFilepath)
	assert.Equal(t, outputLock, lock)
}

func TestGetMainFileFromProjectOutsideRepository(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	relFilepath := "../outside.txt"

	resultFilepath, outputLock, err := postService.GetMainFileFromProject(postID, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", resultFilepath)
	assert.Equal(t, lock, outputLock)
}

func TestGetMainFileFromProjectFailedGetPost(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	relFilepath := "test.txt"

	mockPostRepository.EXPECT().GetByID(postID).Return(nil, errors.New("failed"))

	resultFilepath, outputLock, err := postService.GetMainFileFromProject(postID, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", resultFilepath)
	assert.Equal(t, outputLock, lock)
}

func TestGetMainFileFromProjectFailedCheckoutBranch(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	relFilepath := "test.txt"

	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().LockDirectory(postID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(errors.New("failed"))

	resultFilepath, outputLock, err := postService.GetMainFileFromProject(postID, relFilepath)
	assert.NotNil(t, err)
	assert.Equal(t, "", resultFilepath)
	assert.Equal(t, outputLock, lock)
}

func TestGetMainFileFromProjectFileNotExist(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)
	relFilepath := "child_dir/notreal.txt"

	mockPostRepository.EXPECT().GetByID(postID).Return(&models.Post{Model: gorm.Model{ID: postID}}, nil)
	mockFilesystem.EXPECT().LockDirectory(postID).Return(lock, nil)
	mockFilesystem.EXPECT().CheckoutDirectory(postID)
	mockFilesystem.EXPECT().CheckoutBranch("master").Return(nil)
	mockFilesystem.EXPECT().GetCurrentQuartoDirPath().Return("../utils/test_files/file_tree")

	_, _, err := postService.GetMainFileFromProject(postID, relFilepath)
	assert.NotNil(t, err)
}

func TestGetPost(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	databasePost := &models.Post{
		Model:         gorm.Model{ID: 5},
		Collaborators: []*models.PostCollaborator{},
		Title:         "Hello, world!",
		PostType:      models.Project,
		ScientificFieldTagContainer: models.ScientificFieldTagContainer{
			ScientificFieldTags: []*models.ScientificFieldTag{},
		},
		DiscussionContainer: models.DiscussionContainer{
			Model:       gorm.Model{ID: 6},
			Discussions: []*models.Discussion{},
		},
		DiscussionContainerID: 6,
	}

	mockPostRepository.EXPECT().GetByID(uint(10)).Return(databasePost, nil).Times(1)
	mockPostRepository.EXPECT().GetByID(uint(10)).Return(databasePost, nil).Times(1)

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

	// For this test, we want all posts!
	form := forms.PostFilterForm{
		IncludeProjectPosts: true,
	}

	// Setup mock function return values
	mockPostRepository.EXPECT().QueryPaginated(page, size, gomock.Any()).Return([]*models.Post{
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

	mockPostRepository.EXPECT().QueryPaginated(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("oh no")).Times(1)
	mockPostRepository.EXPECT().QueryPaginated(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := postService.Filter(1, 10, forms.PostFilterForm{})

	if err == nil {
		t.Fatal("post filtering should have failed")
	}
}

func TestGetProjectPostGoodWeather(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)

	post := &models.Post{
		Model: gorm.Model{ID: postID},
	}

	projectPost := &models.ProjectPost{
		Post:   *post,
		PostID: postID,
	}

	// Set up mock function return values
	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil).Times(1)
	mockProjectPostRepository.EXPECT().Query(gomock.Any()).Return([]*models.ProjectPost{projectPost}, nil).Times(1)

	// Function under test
	fetchedProjectPost, err := postService.GetProjectPost(postID)
	if err != nil {
		t.Fatal(err)
	}

	expectedProjectPost := &models.ProjectPost{
		Post: models.Post{
			Model: gorm.Model{ID: postID},
		},
		PostID: postID,
	}

	assert.Equal(t, expectedProjectPost, fetchedProjectPost)
}

func TestGetProjectPostPostNotFound(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)

	mockPostRepository.EXPECT().GetByID(postID).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := postService.GetProjectPost(postID)

	assert.NotNil(t, err)
}

func TestGetProjectPostProjectPostNotFound(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)

	post := &models.Post{
		Model: gorm.Model{ID: postID},
	}

	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil).Times(1)
	mockProjectPostRepository.EXPECT().Query(gomock.Any()).Return([]*models.ProjectPost{}, nil).Times(1)

	// Function under test
	_, err := postService.GetProjectPost(postID)

	assert.NotNil(t, err)
}

func TestGetProjectPostMultipleProjectPostsFound(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)

	post := &models.Post{
		Model: gorm.Model{ID: postID},
	}

	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil).Times(1)

	projectPostA := &models.ProjectPost{
		Post:   models.Post{},
		PostID: postID,
	}

	projectPostB := &models.ProjectPost{
		Post:   models.Post{},
		PostID: postID,
	}

	mockProjectPostRepository.EXPECT().Query(gomock.Any()).Return([]*models.ProjectPost{
		projectPostA,
		projectPostB,
	}, nil).Times(1)

	// Function under test
	_, err := postService.GetProjectPost(postID)

	assert.NotNil(t, err)
}

func TestGetProjectPostProjectPostQueryFailed(t *testing.T) {
	postServiceSetup(t)
	t.Cleanup(postServiceTeardown)

	postID := uint(10)

	post := &models.Post{
		Model: gorm.Model{ID: postID},
	}

	mockPostRepository.EXPECT().GetByID(postID).Return(post, nil).Times(1)
	mockProjectPostRepository.EXPECT().Query(gomock.Any()).Return(nil, fmt.Errorf("oh no")).Times(1)

	// Function under test
	_, err := postService.GetProjectPost(postID)

	assert.NotNil(t, err)
}
