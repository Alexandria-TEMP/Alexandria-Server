package services

import (
	"strconv"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type TagService struct {
	TagRepository database.RepositoryInterface[*tags.ScientificFieldTag]
}

func (tagService *TagService) GetTagsFromStringIDs(ids []string) ([]*tags.ScientificFieldTag, error) {
	tagPointers := []*tags.ScientificFieldTag{}

	for _, s := range ids {
		stringID, err := strconv.ParseUint(s, 10, 64)

		if err != nil {
			return nil, err
		}

		tagID := uint(stringID)

		tag, err := tagService.TagRepository.GetByID(tagID)

		if err != nil {
			return nil, err
		}

		tagPointers = append(tagPointers, tag)
	}

	return tagPointers, nil
}

func (tagService *TagService) GetTagsFromUintIDs(ids []uint) ([]*tags.ScientificFieldTag, error) {
	tagPointers := []*tags.ScientificFieldTag{}

	for _, id := range ids {
		tagID := uint(id)

		tag, err := tagService.TagRepository.GetByID(tagID)

		if err != nil {
			return nil, err
		}

		tagPointers = append(tagPointers, tag)
	}

	return tagPointers, nil
}
