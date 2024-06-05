package models

import (
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestBranchReviewJSONMarshaling(t *testing.T) {
	// This model...
	model := BranchReview{
		Model:                gorm.Model{ID: 88},
		BranchID:             40,
		Member:               Member{},
		MemberID:             50,
		BranchReviewDecision: Approved,
		Feedback:             "Nice!",
	}

	// should equal this DTO!
	targetDTO := BranchReviewDTO{
		ID:                   88,
		BranchID:             40,
		MemberID:             50,
		BranchReviewDecision: Approved,
		Feedback:             "Nice!",
	}

	dto := BranchReviewDTO{}

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
