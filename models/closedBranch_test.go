package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestClosedBranchJSONMarshaling(t *testing.T) {
	supercededBranchID := uint(12)
	// This model...
	model := ClosedBranch{
		Model:              gorm.Model{ID: 55},
		Branch:             Branch{},
		BranchID:           33,
		SupercededBranch:   &Branch{},
		SupercededBranchID: &supercededBranchID,
		ProjectPostID:      40,
		BranchDecision:     Rejected,
	}

	// should equal this DTO!
	targetDTO := ClosedBranchDTO{
		ID:                 55,
		BranchID:           33,
		SupercededBranchID: &supercededBranchID,
		ProjectPostID:      40,
		BranchDecision:     Rejected,
	}

	dto := ClosedBranchDTO{}

	bytes, err := model.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, targetDTO, dto)
}
