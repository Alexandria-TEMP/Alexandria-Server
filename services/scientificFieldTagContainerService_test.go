package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var scientificFieldTagContainerService *ScientificFieldTagContainerService

func setupScientificFieldTagContainer(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockScientificFieldTagContainerRepository = mocks.NewMockModelRepositoryInterface[*models.ScientificFieldTagContainer](mockCtrl)

	// Setup SUT
	scientificFieldTagContainerService = &ScientificFieldTagContainerService{
		ContainerRepository: mockScientificFieldTagContainerRepository,
	}
}

func teardownScientificFieldTagContainer() {

}

func TestGetScientificFieldTagContainer(t *testing.T) {
	setupScientificFieldTagContainer(t)
	t.Cleanup(teardownScientificFieldTagContainer)

	containerID := uint(10)

	databaseContainer := &models.ScientificFieldTagContainer{
		Model: gorm.Model{ID: 2},
		ScientificFieldTags: []*models.ScientificFieldTag{
			{
				Model:           gorm.Model{ID: 5},
				ScientificField: "Mathematics",
			},
			{
				Model:           gorm.Model{ID: 7},
				ScientificField: "Computer Science",
			},
		},
	}

	mockScientificFieldTagContainerRepository.EXPECT().GetByID(containerID).Return(databaseContainer, nil).Times(1)

	// Function under test
	fetchedContainer, err := scientificFieldTagContainerService.GetScientificFieldTagContainer(containerID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, databaseContainer, fetchedContainer)
}
