package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBranchReviewDecisionIsValid(t *testing.T) {
	var good, bad BranchReviewDecision = Approved, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestBranchReviewJSONMarshaling(t *testing.T) {
	createdAt := time.Now().UTC()

	branchReview := &BranchReview{
		Model:                gorm.Model{ID: 5, CreatedAt: createdAt},
		BranchID:             10,
		Member:               Member{Model: gorm.Model{ID: 99}},
		MemberID:             99,
		BranchReviewDecision: Approved,
		Feedback:             "hey nice job",
	}

	// Marshaling the model should convert it to DTO form
	bytes, err := json.Marshal(branchReview)
	if err != nil {
		t.Fatal(err)
	}

	// So unmarshaling, should give us the DTO
	actualDTO := &BranchReviewDTO{}
	if err := json.Unmarshal(bytes, actualDTO); err != nil {
		t.Fatal(err)
	}

	expectedDTO := &BranchReviewDTO{
		ID:                   5,
		BranchID:             10,
		MemberID:             99,
		BranchReviewDecision: Approved,
		Feedback:             "hey nice job",
		CreatedAt:            createdAt,
	}

	assert.Equal(t, expectedDTO, actualDTO)
}
