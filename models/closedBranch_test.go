package models

import (
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestClosedBranchJSONMarshaling(t *testing.T) {
	// This model...
	model := ClosedBranch{
		Model:                gorm.Model{ID: 55},
		Branch:               Branch{},
		BranchID:             33,
		SupercededBranch:     Branch{},
		SupercededBranchID:   12,
		ProjectPostID:        40,
		BranchReviewDecision: Rejected,
	}

	// should equal this DTO!
	targetDTO := ClosedBranchDTO{
		ID:                   55,
		BranchID:             33,
		SupercededBranchID:   12,
		ProjectPostID:        40,
		BranchReviewDecision: Rejected,
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

	if dto != targetDTO {
		t.Fatal("parsed DTO did not equal target DTO")
	}
}
