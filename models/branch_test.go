package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBranchJSONMarshaling(t *testing.T) {
	// This model...
	model := Branch{
		Model: gorm.Model{ID: 44},
		Collaborators: []*BranchCollaborator{
			{Model: gorm.Model{ID: 100}},
			{Model: gorm.Model{ID: 50}},
		},
		Reviews: []*BranchReview{
			{
				Model:                gorm.Model{ID: 2},
				BranchID:             44,
				Member:               Member{},
				MemberID:             88,
				BranchReviewDecision: Approved,
				Feedback:             "LGTM",
			},
		},
		ProjectPostID:             45,
		BranchTitle:               "My Cool MR",
		UpdatedPostTitle:          "Updated Post Title",
		UpdatedCompletionStatus:   Idea,
		UpdatedScientificFields:   []ScientificField{Mathematics},
		RenderStatus:              Pending,
		BranchOverallReviewStatus: BranchOpenForReview,
	}

	// should equal this DTO!
	targetDTO := BranchDTO{
		ID:                        44,
		CollaboratorIDs:           []uint{100, 50},
		ReviewIDs:                 []uint{2},
		ProjectPostID:             45,
		BranchTitle:               "My Cool MR",
		UpdatedPostTitle:          "Updated Post Title",
		UpdatedCompletionStatus:   Idea,
		UpdatedScientificFields:   []ScientificField{Mathematics},
		DiscussionIDs:             []uint{},
		RenderStatus:              Pending,
		BranchOverallReviewStatus: BranchOpenForReview,
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

	assert.Equal(t, targetDTO, dto)
}
