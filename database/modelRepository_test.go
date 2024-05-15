package database

import (
	"log"
	"os"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var modelRepository ModelRepository[*models.Member]
var member models.Member

func setup() {
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

	modelRepository = ModelRepository[*models.Member]{database: testDB}
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	os.Exit(code)
}

// Helper function that deletes database contents
func cleanDatabase() {
	// Delete all members
	testDB.Unscoped().Where("id >= 0").Delete(&models.Member{})
}

func TestCreateWithoutSpecifyingID(t *testing.T) {
	t.Cleanup(cleanDatabase)

	err := modelRepository.Create(&member)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateWithID(t *testing.T) {
	t.Cleanup(cleanDatabase)

	var id uint = 5

	memberWithID := models.Member{
		Model:       gorm.Model{ID: id},
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	}

	err := modelRepository.Create(&memberWithID)
	if err != nil {
		t.Fatal(err)
	}

	if memberWithID.ID != id {
		t.Fatalf("creation did not use ID %d", id)
	}
}

func TestGetById(t *testing.T) {
	t.Cleanup(cleanDatabase)

	// Create a member
	model := models.Member{
		Model:       gorm.Model{ID: 5},
		FirstName:   "first name",
		LastName:    "last name",
		Email:       "email",
		Password:    "password",
		Institution: "institution",
	}

	err := modelRepository.Create(&model)
	if err != nil {
		t.Fatal(err)
	}

	// Try to fetch the member
	id := model.ID
	found, err := modelRepository.GetByID(id)

	if err != nil {
		t.Fatal(err)
	}

	if found.ID != id {
		t.Fatal("fetched ID is not equal to ID at creation time")
	}
}
