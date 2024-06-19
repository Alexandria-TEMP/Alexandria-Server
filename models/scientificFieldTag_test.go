package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestScientificFieldTagJSONMarshaling(t *testing.T) {
	scientificFieldTagID := uint(10)

	scientificFieldTag := &ScientificFieldTag{
		Model:           gorm.Model{ID: scientificFieldTagID},
		ScientificField: "mathematics",
		Containers:      []*ScientificFieldTagContainer{},
		Subtags: []*ScientificFieldTag{
			{
				Model:           gorm.Model{ID: 45},
				ScientificField: "linear algebra",
				Containers:      []*ScientificFieldTagContainer{},
				Subtags:         []*ScientificFieldTag{},
				ParentID:        &scientificFieldTagID,
			},
		},
		ParentID: nil,
	}

	// Marshaling the model should convert it to DTO form
	bytes, err := json.Marshal(scientificFieldTag)
	if err != nil {
		t.Fatal(err)
	}

	// So unmarshaling, should give us the DTO
	actualDTO := &ScientificFieldTagDTO{}
	if err := json.Unmarshal(bytes, actualDTO); err != nil {
		t.Fatal(err)
	}

	expectedDTO := &ScientificFieldTagDTO{
		ID:              scientificFieldTagID,
		ScientificField: "mathematics",
		SubtagIDs:       []uint{45},
		ParentID:        nil,
	}

	assert.Equal(t, expectedDTO, actualDTO)
}
