package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

func TestBranchJSONMarshaling(t *testing.T) {
	// This model...
	model := Branch{
		Model:        gorm.Model{ID: 44},
		NewVersion:   Version{},
		NewVersionID: 99,
		Collaborators: []*BranchCollaborator{
			{Model: gorm.Model{ID: 100}},
			{Model: gorm.Model{ID: 50}},
		},
		Reviews: []*BranchReview{
			{
				Model:          gorm.Model{ID: 2},
				BranchID:       44,
				Member:         Member{},
				MemberID:       88,
				BranchDecision: Approved,
				Feedback:       "LGTM",
			},
		},
		ProjectPostID:           45,
		PreviousVersion:         Version{},
		PreviousVersionID:       20,
		BranchTitle:             "My Cool MR",
		NewPostTitle:            "Updated Post Title",
		UpdatedCompletionStatus: tags.Idea,
		UpdatedScientificFields: []tags.ScientificField{tags.Mathematics},
		Anonymous:               false,
	}

	// should equal this DTO!
	targetDTO := BranchDTO{
		ID:                      44,
		NewVersionID:            99,
		CollaboratorIDs:         []uint{100, 50},
		ReviewIDs:               []uint{2},
		ProjectPostID:           45,
		PreviousVersionID:       20,
		BranchTitle:             "My Cool MR",
		NewPostTitle:            "Updated Post Title",
		UpdatedCompletionStatus: tags.Idea,
		UpdatedScientificFields: []tags.ScientificField{tags.Mathematics},
		Anonymous:               false,
	}

	dto := BranchDTO{}

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
