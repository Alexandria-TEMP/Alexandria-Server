package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRenderStatusIsValid(t *testing.T) {
	var good, bad RenderStatus = Pending, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestBranchOverallReviewStatusIsValid(t *testing.T) {
	var good, bad BranchOverallReviewStatus = BranchOpenForReview, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestBranchJSONMarshaling(t *testing.T) {
	projectPostID := uint(45)
	updatedPostTitle := "Updated Post Title"
	updatedCompletionStatus := Idea

	createdAt := time.Now().Add(time.Minute).UTC()
	updatedAt := time.Now().Add(time.Hour).UTC()
	updatedScientificFieldTagContainerID := uint(20)

	// This model...
	model := Branch{
		Model: gorm.Model{ID: 44, CreatedAt: createdAt, UpdatedAt: updatedAt},
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
		DiscussionContainerID:                5,
		ProjectPostID:                        &projectPostID,
		BranchTitle:                          "My Cool MR",
		UpdatedPostTitle:                     &updatedPostTitle,
		UpdatedCompletionStatus:              &updatedCompletionStatus,
		UpdatedScientificFieldTagContainer:   &ScientificFieldTagContainer{},
		UpdatedScientificFieldTagContainerID: &updatedScientificFieldTagContainerID,
		RenderStatus:                         Pending,
		BranchOverallReviewStatus:            BranchOpenForReview,
	}

	// should equal this DTO!
	targetDTO := BranchDTO{
		ID:                                   44,
		CollaboratorIDs:                      []uint{100, 50},
		ReviewIDs:                            []uint{2},
		ProjectPostID:                        &projectPostID,
		BranchTitle:                          "My Cool MR",
		UpdatedPostTitle:                     &updatedPostTitle,
		UpdatedCompletionStatus:              &updatedCompletionStatus,
		UpdatedScientificFieldTagContainerID: &updatedScientificFieldTagContainerID,
		DiscussionContainerID:                5,
		RenderStatus:                         Pending,
		BranchOverallReviewStatus:            BranchOpenForReview,
		CreatedAt:                            createdAt,
		UpdatedAt:                            updatedAt,
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
