package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestMemberJSONMarshaling(t *testing.T) {
	// This model...
	model := Member{
		Model:                         gorm.Model{ID: 100},
		FirstName:                     "first name",
		LastName:                      "last name",
		Email:                         "email",
		Password:                      "password",
		Institution:                   "institution",
		ScientificFieldTagContainer:   ScientificFieldTagContainer{},
		ScientificFieldTagContainerID: 20,
	}

	// should equal this DTO!
	targetDTO := MemberDTO{
		ID:                            100,
		FirstName:                     "first name",
		LastName:                      "last name",
		Email:                         "email",
		Password:                      "password",
		Institution:                   "institution",
		ScientificFieldTagContainerID: 20,
	}

	dto := MemberDTO{}

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
