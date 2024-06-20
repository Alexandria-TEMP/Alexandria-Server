package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostTypeIsValid(t *testing.T) {
	var good, bad PostType = Question, ""

	assert.True(t, good.IsValid())
	assert.False(t, bad.IsValid())
}

func TestPostJSONMarshaling(t *testing.T) {
	createdAt := time.Now().Add(time.Minute).UTC()
	updatedAt := time.Now().Add(time.Hour).UTC()

	// This model...
	model := Post{
		Model: gorm.Model{ID: 88, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Collaborators: []*PostCollaborator{
			{
				Model:             gorm.Model{ID: 1},
				Member:            Member{},
				MemberID:          90,
				PostID:            88,
				CollaborationType: Author,
			},
			{
				Model:             gorm.Model{ID: 60},
				Member:            Member{},
				MemberID:          20,
				PostID:            88,
				CollaborationType: Contributor,
			},
		},
		DiscussionContainerID:         5,
		Title:                         "Nice Post",
		PostType:                      Question,
		ScientificFieldTagContainer:   ScientificFieldTagContainer{},
		ScientificFieldTagContainerID: 50,
		DiscussionContainer:           DiscussionContainer{Discussions: []*Discussion{{Model: gorm.Model{ID: 95}}}},
	}

	// should equal this DTO!
	targetDTO := PostDTO{
		ID:                            88,
		CollaboratorIDs:               []uint{1, 60},
		Title:                         "Nice Post",
		PostType:                      Question,
		ScientificFieldTagContainerID: 50,
		DiscussionContainerID:         5,
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}

	dto := PostDTO{}

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
