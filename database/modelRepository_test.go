package database

import (
	"log"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var memberRepository ModelRepository[*models.Member]
var projectPostRepository ModelRepository[*models.ProjectPost]
var member models.Member

func beforeEach() {
	database, err := InitializeTestDatabase()
	if err != nil {
		log.Fatalf("Could not initialize test database: %s", err)
	}

	testDB = database

	member = models.Member{
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	}

	memberRepository = ModelRepository[*models.Member]{Database: testDB}
	projectPostRepository = ModelRepository[*models.ProjectPost]{Database: testDB}
}

func afterEach() {
	// Delete all models created by tests
	testDB.Unscoped().Where("id >= 0").Delete(&models.Member{})
	testDB.Unscoped().Where("id >= 0").Delete(&models.Post{})
	testDB.Unscoped().Where("id >= 0").Delete(&models.ProjectPost{})
}

func TestCreateWithoutSpecifyingID(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateWithID(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	var id uint = 5

	memberWithID := models.Member{
		Model:       gorm.Model{ID: id},
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	}

	err := memberRepository.Create(&memberWithID)
	if err != nil {
		t.Fatal(err)
	}

	if memberWithID.ID != id {
		t.Fatalf("creation did not use ID %d", id)
	}
}

func TestGetById(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Create a member
	model := models.Member{
		Model:       gorm.Model{ID: 5},
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	}

	err := memberRepository.Create(&model)
	if err != nil {
		t.Fatal(err)
	}

	// Try to fetch the member
	id := model.ID
	found, err := memberRepository.GetByID(id)

	if err != nil {
		t.Fatal(err)
	}

	if found.ID != id {
		t.Fatal("fetched ID is not equal to ID at creation time")
	}
}

func TestGetByIDReturnsError(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Insert a model with a different ID from the one we're getting
	var idA, idB uint = 5, 66

	member.Model = gorm.Model{ID: idA}

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}

	found, err := memberRepository.GetByID(idB)
	if err == nil {
		t.Fatalf("expected not to find model, but found model with ID %d", found.Model.ID)
	}
}

// Test updating a model, by creating a new model instance with the same ID
func TestUpdateWithNewModelWithSameID(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Insert an initial model
	var id uint = 99
	member.Model = gorm.Model{ID: id}

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}

	// Update the model, using the same ID
	newModel := models.Member{
		Model:       gorm.Model{ID: id},
		FirstName:   "updated first name",
		LastName:    "updated last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	}

	updated, err := memberRepository.Update(&newModel)
	if err != nil {
		t.Fatal(err)
	}

	if updated.FirstName != "updated first name" || updated.LastName != "updated last name" {
		t.Fatal("model fields did not update")
	}

	if updated.Email != "email" {
		t.Fatal("model fields did not retain")
	}
}

// Test updating a model, by getting the original instance and making changes to it
func TestUpdateWithModelFetchedFromDB(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Insert an initial model
	var id uint = 55
	member.Model = gorm.Model{ID: id}

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch the model, change data, and update it
	found, err := memberRepository.GetByID(id)
	if err != nil {
		t.Fatal(err)
	}

	found.FirstName = "new value"

	updated, err := memberRepository.Update(found)
	if err != nil {
		t.Fatal(err)
	}

	if updated.FirstName != "new value" {
		t.Fatal("model was not updated")
	}
}

func TestUpdateWithNonExistingID(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Insert an initial model with a different ID
	var idA, idB uint = 100, 500
	member.Model = gorm.Model{ID: idA}

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}

	// Try to update it, but using a different ID
	newModel := models.Member{
		Model:       gorm.Model{ID: idB},
		FirstName:   "A",
		LastName:    "B",
		Email:       "C",
		Password:    "D",
		Institution: "E",
	}

	_, err = memberRepository.Update(&newModel)
	if err == nil {
		t.Fatal("expected error after update using new ID")
	}
}

func TestDeleteExistingModel(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Insert an initial model
	var id uint = 100
	member.Model = gorm.Model{ID: id}

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}

	// Delete it again
	err = memberRepository.Delete(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteNonExistingModel(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Insert an initial model with a different ID
	var idA, idB uint = 100, 500
	member.Model = gorm.Model{ID: idA}

	err := memberRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}

	// Delete a different ID
	err = memberRepository.Delete(idB)
	if err == nil {
		t.Fatal("deletion should have failed")
	}
}

// Project Post contains Post. When we fetch Project Post, we want to pre-load Post
// so that it is also included in the result. We test that here.
func TestGetPreloadedAssociations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	createdProjectPost := models.ProjectPost{
		Post: models.Post{
			Collaborators:       []*models.PostCollaborator{},
			Title:               "TEST POST",
			PostType:            tags.Project,
			ScientificFieldTags: []tags.ScientificField{},
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

	if err := projectPostRepository.Create(&createdProjectPost); err != nil {
		log.Fatal(err)
	}

	fetchedProjectPost, err := projectPostRepository.GetByID(createdProjectPost.ID)
	if err != nil {
		log.Fatal(err)
	}

	// If pre-loading worked, the nested Post's fields will be updated.
	if !(fetchedProjectPost.Post.Title == "TEST POST") {
		t.Fatal("nested Post object did not pre load")
	}
}
