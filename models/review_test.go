package models

import (
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestReviewJSONMarshaling(t *testing.T) {
	// This model...
	model := Review{
		Model:          gorm.Model{ID: 88},
		BranchID:       40,
		Member:         Member{},
		MemberID:       50,
		BranchDecision: Approved,
		Feedback:       "Nice!",
	}

	// should equal this DTO!
	targetDTO := ReviewDTO{
		ID:             88,
		BranchID:       40,
		MemberID:       50,
		BranchDecision: Approved,
		Feedback:       "Nice!",
	}

	dto := ReviewDTO{}

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
