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
		PostType:              Project,
		ScientificFieldTags:   []tags.ScientificField{tags.Mathematics},
		DiscussionContainer:   DiscussionContainer{Model: gorm.Model{ID: 50}, Discussions: []*Discussion{{Model: gorm.Model{ID: 95}}}},
		DiscussionContainerID: 50,
	}

	model := ProjectPost{
		Model:        gorm.Model{ID: 42},
		Post:         post,
		PostID:       88,
		OpenBranches: []*Branch{{Model: gorm.Model{ID: 44}}},
		ClosedBranches: []*ClosedBranch{
			{Model: gorm.Model{ID: 59}},
			{Model: gorm.Model{ID: 20}},
		},
		CompletionStatus:   Completed,
		FeedbackPreference: FormalFeedback,
		PostReviewStatus:   RevisionNeeded,
	}

	// should equal this DTO!
	targetDTO := ProjectPostDTO{
		ID: 42,
		PostDTO: PostDTO{
			ID:                    88,
			CollaboratorIDs:       []uint{1, 60},
			PostType:              Project,
			ScientificFieldTags:   []tags.ScientificField{tags.Mathematics},
			DiscussionContainerID: 50,
		},
		OpenBranchIDs:      []uint{44},
		ClosedBranchIDs:    []uint{59, 20},
		CompletionStatus:   Completed,
		FeedbackPreference: FormalFeedback,
		PostReviewStatus:   RevisionNeeded,
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
		t.Fatalf("parsed DTO\n%+v\ndid not equal target DTO\n%+v", dto, targetDTO)
	}
}
