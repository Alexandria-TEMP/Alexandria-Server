package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type TagService struct {
}

func (tagService *TagService) GetTagsFromIDs(_ []string) ([]*tags.ScientificFieldTag, error) {
	tagPointers := []*tags.ScientificFieldTag{}
	var err error
	err = nil
	return tagPointers, err
}
