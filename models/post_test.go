package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
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
		CurrentVersion:      Version{},
		CurrentVersionID:    49,
		PostType:            tags.Question,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
	}

	// should equal this DTO!
	targetDTO := PostDTO{
		ID:                  88,
		CollaboratorIDs:     []uint{1, 60},
		VersionID:           49,
		PostType:            tags.Question,
		ScientificFieldTags: []tags.ScientificField{tags.Mathematics},
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
