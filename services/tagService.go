package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type TagService struct {
	TagRepository database.ModelRepositoryInterface[*tags.ScientificFieldTag]
	// TagContainerRepository database.RepositoryInterface[*tags.ScientificFieldTagContainer]
}

func (tagService *TagService) GetAllScientificFieldTags() ([]*tags.ScientificFieldTag, error) {
	tags, err := tagService.TagRepository.Query()
	return tags, err
}

func (tagService *TagService) GetTagsFromUintIDs(ids []uint) ([]*tags.ScientificFieldTag, error) {
	tagPointers := []*tags.ScientificFieldTag{}

	for _, id := range ids {
		tag, err := tagService.TagRepository.GetByID(id)

		if err != nil {
			return nil, err
		}

		tagPointers = append(tagPointers, tag)
	}

	return tagPointers, nil
}
