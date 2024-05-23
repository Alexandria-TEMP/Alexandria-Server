package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

func TestMergeRequestJSONMarshaling(t *testing.T) {
	// This model...
	model := MergeRequest{
		Model:        gorm.Model{ID: 44},
		NewVersion:   Version{},
		NewVersionID: 99,
		Collaborators: []*MergeRequestCollaborator{
			{
				Model: gorm.Model{ID: 100},
			},
			{
				Model: gorm.Model{ID: 50},
			},
		},
		Reviews: []*MergeRequestReview{
			{
				Model:                gorm.Model{ID: 2},
				MergeRequestID:       44,
				Member:               Member{},
				MemberID:             88,
				MergeRequestDecision: Approved,
				Feedback:             "LGTM",
			},
		},
		ProjectPostID:           45,
		Title:                   "My Cool MR",
		UpdatedCompletionStatus: tags.Idea,
		UpdatedScientificFields: tags.Mathematics,
		Anonymous:               false,
	}

	// should equal this DTO!
	targetDTO := MergeRequestDTO{
		ID:                      44,
		NewVersionID:            99,
		CollaboratorIDs:         []uint{100, 50},
		ReviewIDs:               []uint{2},
		ProjectPostID:           45,
		Title:                   "My Cool MR",
		UpdatedCompletionStatus: tags.Idea,
		UpdatedScientificFields: tags.Mathematics,
		Anonymous:               false,
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
