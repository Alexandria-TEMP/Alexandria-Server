package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProjectCompletionStatusIsValid(t *testing.T) {
	var good, bad ProjectCompletionStatus = Idea, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestProjectFeedbackPreference(t *testing.T) {
	var good, bad ProjectFeedbackPreference = DiscussionFeedback, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestProjectReviewStatusIsValid(t *testing.T) {
	var good, bad ProjectReviewStatus = RevisionNeeded, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestProjectPostJSONMarshaling(t *testing.T) {
	// This model...
	post := Post{
		Model: gorm.Model{ID: 88},
		Collaborators: []*PostCollaborator{
			{Model: gorm.Model{ID: 1}},
			{Model: gorm.Model{ID: 60}},
		},
		PostType:                    Project,
		ScientificFieldTagContainer: ScientificFieldTagContainer{},
		DiscussionContainer:         DiscussionContainer{Model: gorm.Model{ID: 50}, Discussions: []*Discussion{{Model: gorm.Model{ID: 95}}}},
		DiscussionContainerID:       50,
	}

	createdAt := time.Now().Add(time.Minute).UTC()
	updatedAt := time.Now().Add(time.Hour).UTC()

	model := ProjectPost{
		Model:        gorm.Model{ID: 42, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Post:         post,
		PostID:       88,
		OpenBranches: []*Branch{{Model: gorm.Model{ID: 44}}},
		ClosedBranches: []*ClosedBranch{
			{Model: gorm.Model{ID: 59}},
			{Model: gorm.Model{ID: 20}},
		},
		ProjectCompletionStatus:   Completed,
		ProjectFeedbackPreference: FormalFeedback,
		PostReviewStatus:          RevisionNeeded,
	}

	// should equal this DTO!
	targetDTO := ProjectPostDTO{
		ID:                        42,
		PostID:                    88,
		OpenBranchIDs:             []uint{44},
		ClosedBranchIDs:           []uint{59, 20},
		ProjectCompletionStatus:   Completed,
		ProjectFeedbackPreference: FormalFeedback,
		PostReviewStatus:          RevisionNeeded,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
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
