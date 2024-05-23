package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

func TestProjectPostJSONMarshaling(t *testing.T) {
	// This model...
	post := Post{
		Model: gorm.Model{ID: 88},
		Collaborators: []*PostCollaborator{
			{Model: gorm.Model{ID: 1}},
			{Model: gorm.Model{ID: 60}},
		},
		CurrentVersion:      Version{},
		CurrentVersionID:    49,
		PostType:            tags.Project,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}

	model := ProjectPost{
		Model:  gorm.Model{ID: 42},
		Post:   post,
		PostID: 88,
		OpenMergeRequests: []*MergeRequest{
			{Model: gorm.Model{ID: 44}},
		},
		ClosedMergeRequests: []*ClosedMergeRequest{
			{Model: gorm.Model{ID: 59}},
			{Model: gorm.Model{ID: 20}},
		},
		HasHadInitialPeerReview: true,
		CompletionStatus:        tags.Completed,
		FeedbackPreference:      tags.FormalFeedback,
		PostReviewStatusTag:     tags.RevisionNeeded,
	}

	// should equal this DTO!
	targetDTO := ProjectPostDTO{
		ID: 42,
		PostDTO: PostDTO{
			ID:                  88,
			CollaboratorIDs:     []uint{1, 60},
			VersionID:           49,
			PostType:            tags.Project,
			ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
		},
		HasHadInitialPeerReview: true,
		OpenMergeRequestIDs:     []uint{44},
		ClosedMergeRequestIDs:   []uint{59, 20},
		CompletionStatus:        tags.Completed,
		FeedbackPreference:      tags.FormalFeedback,
		PostReviewStatusTag:     tags.RevisionNeeded,
	}

	dto := ProjectPostDTO{}

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
