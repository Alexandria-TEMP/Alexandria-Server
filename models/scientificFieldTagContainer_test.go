package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestScientificFieldTagContainerJSONMarshaling(t *testing.T) {
	scientificFieldTagContainer := &ScientificFieldTagContainer{
		Model: gorm.Model{ID: 5},
		ScientificFieldTags: []*ScientificFieldTag{
			{
				Model:           gorm.Model{ID: 2},
				ScientificField: "Mathematics",
			},
			{
				Model:           gorm.Model{ID: 3},
				ScientificField: "Computer Science",
			},
		},
	}

	bytes, err := json.Marshal(scientificFieldTagContainer)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %s", err)
	}

	var fetchedScientificFieldTagContainerDTO ScientificFieldTagContainerDTO
	if err := json.Unmarshal(bytes, &fetchedScientificFieldTagContainerDTO); err != nil {
		t.Fatalf("failed to unmarshal JSON: %s", err)
	}

	expectedScientificFieldTagContainerDTO := ScientificFieldTagContainerDTO{
		ID:                    5,
		ScientificFieldTagIDs: []uint{2, 3},
	}

	assert.Equal(t, expectedScientificFieldTagContainerDTO, fetchedScientificFieldTagContainerDTO)
}
