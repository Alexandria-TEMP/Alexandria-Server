package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

func createTestMR(createdAt, updatedAt time.Time) MergeRequest {
	return MergeRequest{
		Model:        gorm.Model{ID: 44, CreatedAt: createdAt, UpdatedAt: updatedAt},
		NewVersion:   Version{},
		NewVersionID: 99,
		Collaborators: []*MergeRequestCollaborator{
			{Model: gorm.Model{ID: 100}},
			{Model: gorm.Model{ID: 50}},
		},
		Reviews: []*MergeRequestReview{
			{
				Model:                gorm.Model{ID: 2},
				MergeRequestID:       44,
				Member:               Member{},
				MemberID:             88,
				MergeRequestDecision: ReviewApproved,
				Feedback:             "LGTM",
			},
		},
		ProjectPostID:           45,
		PreviousVersion:         Version{},
		PreviousVersionID:       20,
		MergeRequestTitle:       "My Cool MR",
		UpdatedPostTitle:        "Updated Post Title",
		UpdatedCompletionStatus: tags.Idea,
		UpdatedScientificFields: []*tags.ScientificFieldTag{},
		Anonymous:               false,
		MergeRequestDecision:    MergeRequestOpenForReview,
	}
}

func TestMergeRequestJSONMarshaling(t *testing.T) {
	createdAt := time.Now().UTC()
	updatedAt := time.Now().Add(time.Hour).UTC()

	// This model...
	model := createTestMR(createdAt, updatedAt)

	// should equal this DTO!
	targetDTO := MergeRequestDTO{
		ID:                           44,
		NewVersionID:                 99,
		CollaboratorIDs:              []uint{100, 50},
		ReviewIDs:                    []uint{2},
		ProjectPostID:                45,
		PreviousVersionID:            20,
		MergeRequestTitle:            "My Cool MR",
		NewPostTitle:                 "Updated Post Title",
		UpdatedCompletionStatus:      tags.Idea,
		UpdatedScientificFieldTagIDs: []uint{},
		Anonymous:                    false,
		MergeRequestDecision:         MergeRequestOpenForReview,
		CreatedAt:                    createdAt,
		UpdatedAt:                    updatedAt,
	}

	dto := MergeRequestDTO{}

	bytes, err := model.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dto, targetDTO) {
		t.Fatal("parsed DTO did not equal target DTO")
	}
}
