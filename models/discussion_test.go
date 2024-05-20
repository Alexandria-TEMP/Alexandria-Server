package models

import (
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestDiscussionJSONMarshaling(t *testing.T) {
	var parentID uint = 5

	// This model...
	model := Discussion{
		Model:     gorm.Model{ID: 100},
		VersionID: 33,
		Member:    Member{},
		MemberID:  28,
		Replies:   []*Discussion{},
		ParentID:  &parentID,
		Text:      "Test!",
		Deleted:   false,
		Anonymous: true,
	}

	// should equal this DTO!
	targetDTO := DiscussionDTO{
		ID:        100,
		VersionID: 33,
		MemberID:  28,
		ParentID:  5,
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

	if dto != targetDTO {
		t.Fatal("parsed DTO did not equal target DTO")
	}
}
