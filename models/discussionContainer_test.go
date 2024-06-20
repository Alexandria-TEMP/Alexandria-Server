package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestDiscussionContainerJSONMarshaling(t *testing.T) {
	rootDiscussionID := uint(10)

	discussionContainer := DiscussionContainer{
		Model: gorm.Model{ID: 5},
		Discussions: []*Discussion{
			{
				Model:       gorm.Model{ID: rootDiscussionID},
				ContainerID: 5,
			},
			{
				Model:       gorm.Model{ID: 12},
				ContainerID: 5,
				ParentID:    &rootDiscussionID,
			},
		},
	}

	marshaled, err := json.Marshal(&discussionContainer)
	if err != nil {
		t.Fatal(err)
	}

	var createdDTO DiscussionContainerDTO
	if err := json.Unmarshal(marshaled, &createdDTO); err != nil {
		t.Fatal(err)
	}

	expectedDTO := DiscussionContainerDTO{
		ID:            5,
		DiscussionIDs: []uint{rootDiscussionID},
	}

	if !reflect.DeepEqual(createdDTO, expectedDTO) {
		t.Fatalf("created DTO\n%+v\nshould equal expected DTO\n%+v", createdDTO, expectedDTO)
	}
}
