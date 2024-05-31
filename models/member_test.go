package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gorm.io/gorm"
)

func TestMemberJSONMarshaling(t *testing.T) {
	// This model...
	model := Member{
		Model:               gorm.Model{ID: 100},
		FirstName:           "first name",
		LastName:            "last name",
		Email:               "email",
		Password:            "password",
		Institution:         "institution",
		ScientificFieldTags: []*tags.ScientificFieldTag{},
	}

	// should equal this DTO!
	targetDTO := MemberDTO{
		ID:                    100,
		FirstName:             "first name",
		LastName:              "last name",
		Email:                 "email",
		Password:              "password",
		Institution:           "institution",
		ScientificFieldTagIDs: []uint{},
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
