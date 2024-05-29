package models

import (
	"encoding/json"
	"testing"

	"gorm.io/gorm"
)

func TestPostCollaboratorJSONMarshaling(t *testing.T) {
	// This model...
	model := PostCollaborator{
		Model:             gorm.Model{ID: 55},
		Member:            Member{},
		MemberID:          32,
		PostID:            87,
		CollaborationType: Author,
	}

	// should equal this DTO!
	targetDTO := PostCollaboratorDTO{
		ID:                55,
		MemberID:          32,
		PostID:            87,
		CollaborationType: Author,
	}

	dto := PostCollaboratorDTO{}

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

func TestBranchCollaboratorJSONMarshaling(t *testing.T) {
	// This model...
	model := BranchCollaborator{
		Model:             gorm.Model{ID: 55},
		Member:            Member{},
		MemberID:          32,
		BranchID:          87,
		CollaborationType: Author,
	}

	// should equal this DTO!
	targetDTO := BranchCollaboratorDTO{
		ID:                55,
		MemberID:          32,
		BranchID:          87,
		CollaborationType: Author,
	}

	dto := BranchCollaboratorDTO{}

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
