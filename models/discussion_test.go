package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestDiscussionJSONMarshaling(t *testing.T) {
	var parentID uint = 5

	// This model...
	model := Discussion{
		Model:    gorm.Model{ID: 100},
		Member:   Member{},
		MemberID: 28,
		Replies: []*Discussion{
			{
				Model: gorm.Model{ID: 50},
			},
			{
				Model: gorm.Model{ID: 88},
			},
		},
		ParentID:  parentID,
		Text:      "Test!",
		Deleted:   false,
		Anonymous: true,
	}

	// should equal this DTO!
	targetDTO := DiscussionDTO{
		ID:        100,
		MemberID:  28,
		ReplyIDs:  []uint{50, 88},
		Text:      "Test!",
		Deleted:   false,
		Anonymous: true,
	}

	dto := DiscussionDTO{}

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
