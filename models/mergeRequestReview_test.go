package models

import (
	"encoding/json"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestMergeRequestReviewJSONMarshaling(t *testing.T) {
	createdAt := time.Now().UTC()

	// This model...
	model := MergeRequestReview{
		Model:                gorm.Model{ID: 88, CreatedAt: createdAt},
		MergeRequestID:       40,
		Member:               Member{},
		MemberID:             50,
		MergeRequestDecision: ReviewApproved,
		Feedback:             "Nice!",
	}

	// should equal this DTO!
	targetDTO := MergeRequestReviewDTO{
		ID:                   88,
		MergeRequestID:       40,
		MemberID:             50,
		MergeRequestDecision: ReviewApproved,
		Feedback:             "Nice!",
		CreatedAt:            createdAt,
	}

	dto := MergeRequestReviewDTO{}

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
