package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestVersionJSONMarshaling(t *testing.T) {
	// This model...
	model := Version{
		Model: gorm.Model{ID: 20},
		Discussions: []*Discussion{
			{
				Model:     gorm.Model{ID: 99},
				VersionID: 20,
			},
			{
				Model:     gorm.Model{ID: 59},
				VersionID: 20,
			},
		},
		RenderStatus: RenderPending,
	}

	// should equal this DTO!
	targetDTO := VersionDTO{
		ID:            20,
		DiscussionIDs: []uint{99, 59},
		RenderStatus:  RenderPending,
	}

	dto := VersionDTO{}

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
