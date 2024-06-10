package database

import (
	"log"
	"reflect"
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
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
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
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
	}

	err := memberRepository.Create(&memberWithID)
	if err != nil {
		t.Fatal(err)
	}

	if memberWithID.ID != id {
		t.Fatalf("creation did not use ID %d", id)
	}
}

func TestCreateFails(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	memberA := models.Member{
		Model: gorm.Model{ID: 5},
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
	}

	err := memberRepository.Create(&memberA)

	if err != nil {
		t.Fatalf("could not create first member: %s", err)
	}

	err = memberRepository.Create(&memberA)

	if err == nil {
		t.Fatal("creation should have returned error")
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
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
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
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
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
		ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
			ScientificFieldTags: []*tags.ScientificFieldTag{},
		},
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
			Collaborators: []*models.PostCollaborator{},
			Title:         "TEST POST",
			PostType:      models.Project,
			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
				ScientificFieldTags: []*tags.ScientificFieldTag{},
			},
			DiscussionContainer: models.DiscussionContainer{
				Discussions: []*models.Discussion{},
			},
		},
		OpenBranches:       []*models.Branch{},
		ClosedBranches:     []*models.ClosedBranch{},
		CompletionStatus:   models.Ongoing,
		FeedbackPreference: models.FormalFeedback,
		PostReviewStatus:   models.Open,
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

// func TestGetQueryFieldsSimple(t *testing.T) {
// 	if testing.Short() {
// 		t.SkipNow()
// 	}

// 	beforeEach()
// 	t.Cleanup(afterEach)
// 	fmt.Println("All is well up to point 1")
// 	// Create dummy members in the database, with specific IDs
// 	membersToCreate := []models.Member{
// 		{Model: gorm.Model{ID: 5},
// 			FirstName:   "one",
// 			LastName:    "One",
// 			Institution: "TU Delft",
// 			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
// 				ScientificFieldTags: []*tags.ScientificFieldTag{},
// 			}},
// 		{Model: gorm.Model{ID: 10},
// 			FirstName:   "two",
// 			LastName:    "Two",
// 			Institution: "Vrije Universiteit Berlin",
// 			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
// 				ScientificFieldTags: []*tags.ScientificFieldTag{},
// 			}},
// 		{Model: gorm.Model{ID: 12},
// 			FirstName:   "three",
// 			LastName:    "Three",
// 			Institution: "Politechnika Poznanska",
// 			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
// 				ScientificFieldTags: []*tags.ScientificFieldTag{},
// 			}},
// 	}

// 	for _, member := range membersToCreate {
// 		if err := memberRepository.Create(&member); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	fmt.Println("All is well up to point 2")

// 	// Try to get all members in the database
// 	fetchedMembers, err := memberRepository.GetFields([]models.MemberShortFormDTO{})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println("All is well up to point 3")
// 	// cast result as the expected data type

// 	memberShortFormDTOs := fetchedMembers.([]models.MemberShortFormDTO)
// 	fmt.Println("All is well up to point 4")
// 	// define expected result
// 	expectedMemberShortFormDTOs := []models.MemberShortFormDTO{
// 		{
// 			ID:        5,
// 			FirstName: "one",
// 			LastName:  "One",
// 		},
// 		{
// 			ID:        10,
// 			FirstName: "two",
// 			LastName:  "Two",
// 		},
// 		{
// 			ID:        12,
// 			FirstName: "three",
// 			LastName:  "Three",
// 		},
// 	}

// 	// Check that the results have the same length
// 	if len(memberShortFormDTOs) != len(expectedMemberShortFormDTOs) {
// 		t.Fatalf("expected %d records, got %d", len(expectedMemberShortFormDTOs), len(memberShortFormDTOs))
// 	}

// 	fmt.Println("All is well up to point 5")

// 	// Check each fetched member's ID, first and last name
// 	for i, fetchedMember := range memberShortFormDTOs {
// 		if fetchedMember.ID != expectedMemberShortFormDTOs[i].ID {
// 			t.Fatalf("encountered ID %d expecting %d", fetchedMember.ID, expectedMemberShortFormDTOs[i].ID)
// 		}
// 		if fetchedMember.FirstName != expectedMemberShortFormDTOs[i].FirstName {
// 			t.Fatalf("encountered first name %s expecting %s", fetchedMember.FirstName, expectedMemberShortFormDTOs[i].FirstName)
// 		}
// 		if fetchedMember.LastName != expectedMemberShortFormDTOs[i].LastName {
// 			t.Fatalf("encountered last name %s expecting %s", fetchedMember.LastName, expectedMemberShortFormDTOs[i].LastName)
// 		}
// 	}
// }

func TestQuerySimple(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Create dummy members in the database, with specific IDs
	membersToCreate := []models.Member{
		{Model: gorm.Model{ID: 5},
			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
				ScientificFieldTags: []*tags.ScientificFieldTag{},
			}},
		{Model: gorm.Model{ID: 10},
			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
				ScientificFieldTags: []*tags.ScientificFieldTag{},
			}},
		{Model: gorm.Model{ID: 12},
			ScientificFieldTagContainer: tags.ScientificFieldTagContainer{
				ScientificFieldTags: []*tags.ScientificFieldTag{},
			}},
	}

	for _, member := range membersToCreate {
		if err := memberRepository.Create(&member); err != nil {
			t.Fatal(err)
		}
	}

	// Try to fetch some via a query
	fetchedMembers, err := memberRepository.Query("id > 6")
	if err != nil {
		t.Fatal(err)
	}

	expectedIDs := []uint{10, 12}

	if len(fetchedMembers) != len(expectedIDs) {
		t.Fatalf("expected %d records, got %d", len(expectedIDs), len(fetchedMembers))
	}

	// Check each fetched member's ID
	for i, fetchedMember := range fetchedMembers {
		if fetchedMember.ID != expectedIDs[i] {
			t.Fatalf("encountered ID %d expecting %d", fetchedMember.ID, expectedIDs[i])
		}
	}
}

func TestQueryNonExistingField(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	// Function under test
	_, err := memberRepository.Query("AsadhjahDJWduadua = DAASDc!121*@@@@")

	if err == nil {
		t.Fatal("nonsense query should have thrown error")
	}

	// Also test paginated here, since it should be exactly the same
	_, err = memberRepository.QueryPaginated(1, 1, "ASDjASHDjahsd123 = ASD2@@@")

	if err == nil {
		t.Fatal("nonsense paginated query should have thrown error")
	}
}

func TestQueryPaginated(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	beforeEach()
	t.Cleanup(afterEach)

	memberIDs := []uint{
		10, 11, 12, 13, 15, 20, 21, 41,
		42, 43, 60, 61, 62, 78, 88,
	}
	container := tags.ScientificFieldTagContainer{
		ScientificFieldTags: []*tags.ScientificFieldTag{},
	}
	// Add members with the above IDs to the database
	for _, memberID := range memberIDs {
		if err := memberRepository.Create(&models.Member{Model: gorm.Model{ID: memberID}, ScientificFieldTagContainer: container}); err != nil {
			t.Fatal(err)
		}
	}

	// Test paginated queries on this created data
	size := 4
	condition := "id >= 20"

	expectedPages := [][]uint{
		{20, 21, 41, 42},
		{43, 60, 61, 62},
		{78, 88},
	}

	// For each page in the expected pages, perform a paginated query,
	// and compare the outcome against expected.
	for i, expectedMemberIDs := range expectedPages {
		// The page number we're querying
		page := i + 1

		fetchedMembers, err := memberRepository.QueryPaginated(page, size, condition)
		if err != nil {
			t.Fatal(err)
		}

		// Extract member IDs for easier comparisons
		fetchedMemberIDs := make([]uint, len(fetchedMembers))
		for i, member := range fetchedMembers {
			fetchedMemberIDs[i] = member.ID
		}

		if !reflect.DeepEqual(fetchedMemberIDs, expectedMemberIDs) {
			t.Fatalf("fetched member IDs\n%+v\ndid not equal expected member IDs\n%+v", fetchedMemberIDs, expectedMemberIDs)
		}
	}
}
