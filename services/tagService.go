package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type TagService struct {
	TagRepository database.ModelRepositoryInterface[*tags.ScientificFieldTag]
	// TagContainerRepository database.RepositoryInterface[*tags.ScientificFieldTagContainer]
}

func (tagService *TagService) GetTagByID(id uint) (*tags.ScientificFieldTag, error) {
	returnedTag, err := tagService.TagRepository.GetByID(id)
	return returnedTag, err
}

func (tagService *TagService) GetAllScientificFieldTags() ([]*tags.ScientificFieldTag, error) {
	returnedTags, err := tagService.TagRepository.Query()
	return returnedTags, err
}

func (tagService *TagService) GetTagsFromIDs(ids []uint) ([]*tags.ScientificFieldTag, error) {
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
