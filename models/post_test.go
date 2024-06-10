package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestPostJSONMarshaling(t *testing.T) {
	// This model...
	model := Post{
		Model: gorm.Model{ID: 88},
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
		Title:               "Nice Post",
		PostType:            Question,
		ScientificFields:    []ScientificField{Mathematics},
		DiscussionContainer: DiscussionContainer{Discussions: []*Discussion{{Model: gorm.Model{ID: 95}}}},
	}

	// should equal this DTO!
	targetDTO := PostDTO{
		ID:               88,
		CollaboratorIDs:  []uint{1, 60},
		Title:            "Nice Post",
		PostType:         Question,
		ScientificFields: []ScientificField{Mathematics},
		DiscussionIDs:    []uint{95},
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
