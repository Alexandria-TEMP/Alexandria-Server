package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestDiscussionJSONMarshaling(t *testing.T) {
	var parentID uint = 5

	var memberID uint = 28

	createdAt := time.Now().UTC()

	// This model...
	model := Discussion{
		Model:    gorm.Model{ID: 100, CreatedAt: createdAt},
		Member:   &Member{},
		MemberID: &memberID,
		Replies: []*Discussion{
			{
				Model: gorm.Model{ID: 50},
			},
			{
				Model: gorm.Model{ID: 88},
			},
		},
		ParentID: &parentID,
		Text:     "Test!",
	}

	// should equal this DTO!
	targetDTO := DiscussionDTO{
		ID:        100,
		MemberID:  &memberID,
		ReplyIDs:  []uint{50, 88},
		Text:      "Test!",
		CreatedAt: createdAt,
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
