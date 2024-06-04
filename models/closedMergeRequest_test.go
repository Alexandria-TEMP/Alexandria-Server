package models

import (
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestClosedMergeRequestJSONMarshaling(t *testing.T) {
	// This model...
	model := ClosedMergeRequest{
		Model:                   gorm.Model{ID: 55},
		MergeRequest:            MergeRequest{},
		MergeRequestID:          33,
		MainVersionWhenClosed:   Version{},
		MainVersionWhenClosedID: 87,
		ProjectPostID:           40,
		MergeRequestDecision:    ReviewRejected,
	}

	// should equal this DTO!
	targetDTO := ClosedMergeRequestDTO{
		ID:                      55,
		MergeRequestID:          33,
		MainVersionWhenClosedID: 87,
		ProjectPostID:           40,
		MergeRequestDecision:    ReviewRejected,
	}

	dto := ClosedMergeRequestDTO{}

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
