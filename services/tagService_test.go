package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// SUT
var tagService TagService

func setupTagService(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Setup mocks
	mockScientificFieldTagRepository = mocks.NewMockModelRepositoryInterface[*models.ScientificFieldTag](mockCtrl)

	// Setup SUT
	tagService = TagService{
		TagRepository: mockScientificFieldTagRepository,
	}
}

func teardownTagService() {

}

func TestGetTagByID(t *testing.T) {
	setupTagService(t)
	t.Cleanup(teardownTagService)

	// Setup data
	tagID := uint(10)

	databaseTag := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: 5},
		ScientificField: "mathematics",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags: []*models.ScientificFieldTag{
			{
				Model:           gorm.Model{ID: 7},
				ScientificField: "linear algebra",
				Containers:      []*models.ScientificFieldTagContainer{},
				Subtags:         []*models.ScientificFieldTag{},
				ParentID:        &tagID,
			},
		},
		ParentID: nil,
	}

	// Setup mocks
	mockScientificFieldTagRepository.EXPECT().GetByID(tagID).Return(databaseTag, nil).Times(1)

	// Function under test
	actualTag, err := tagService.GetTagByID(tagID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, databaseTag, actualTag)
}

func TestGetAllScientificFieldTags(t *testing.T) {
	setupTagService(t)
	t.Cleanup(teardownTagService)

	// Setup data
	parentTagID := uint(10)
	childTagID := uint(5)

	childTag := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: childTagID},
		ScientificField: "linear algebra",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags:         []*models.ScientificFieldTag{},
		ParentID:        &parentTagID,
	}

	parentTag := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: parentTagID},
		ScientificField: "mathematics",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags: []*models.ScientificFieldTag{
			childTag,
		},
		ParentID: nil,
	}

	databaseTags := []*models.ScientificFieldTag{
		childTag,
		parentTag,
	}

	// Setup mocks
	mockScientificFieldTagRepository.EXPECT().Query(gomock.Any()).Return(databaseTags, nil).Times(1)

	// Function under test
	actualTags, err := tagService.GetAllScientificFieldTags()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, databaseTags, actualTags)
}

func TestGetTagsFromIDs(t *testing.T) {
	setupTagService(t)
	t.Cleanup(teardownTagService)

	// Setup data
	parentTagID := uint(10)
	childTagID := uint(5)
	thirdTagID := uint(7)

	childTag := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: childTagID},
		ScientificField: "linear algebra",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags:         []*models.ScientificFieldTag{},
		ParentID:        &parentTagID,
	}

	parentTag := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: parentTagID},
		ScientificField: "mathematics",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags: []*models.ScientificFieldTag{
			childTag,
		},
		ParentID: nil,
	}

	thirdTag := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: thirdTagID},
		ScientificField: "computer science",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags:         []*models.ScientificFieldTag{},
		ParentID:        nil,
	}

	ids := []uint{childTagID, thirdTagID}

	// Setup mocks
	mockScientificFieldTagRepository.EXPECT().GetByID(childTagID).Return(childTag, nil).AnyTimes()
	mockScientificFieldTagRepository.EXPECT().GetByID(parentTagID).Return(parentTag, nil).AnyTimes()
	mockScientificFieldTagRepository.EXPECT().GetByID(thirdTagID).Return(thirdTag, nil).AnyTimes()

	// Function under test
	actualTags, err := tagService.GetTagsFromIDs(ids)
	if err != nil {
		t.Fatal(err)
	}

	expectedTags := []*models.ScientificFieldTag{
		childTag,
		thirdTag,
	}

	assert.Equal(t, expectedTags, actualTags)
}

func TestGetTagsFromIDsTagNotFound(t *testing.T) {
	setupTagService(t)
	t.Cleanup(teardownTagService)

	// Setup data
	idA := uint(3)
	idB := uint(7)

	ids := []uint{idA, idB}

	tagA := &models.ScientificFieldTag{
		Model:           gorm.Model{ID: idA},
		ScientificField: "mathematics",
		Containers:      []*models.ScientificFieldTagContainer{},
		Subtags:         []*models.ScientificFieldTag{},
		ParentID:        nil,
	}

	// Setup mocks
	mockScientificFieldTagRepository.EXPECT().GetByID(idA).Return(tagA, nil)
	mockScientificFieldTagRepository.EXPECT().GetByID(idB).Return(nil, fmt.Errorf("not found"))

	// Function under test
	_, err := tagService.GetTagsFromIDs(ids)

	assert.NotNil(t, err)
}
