package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type TagService struct {
	TagRepository database.RepositoryInterface[*tags.ScientificFieldTag]
	//TagContainerRepository database.RepositoryInterface[*tags.ScientificFieldTagContainer]
}

// func (tagService *TagService) GetTagContainer(tagID uint) (*tags.ScientificFieldTagContainer, error) {
// 	// get Tag by this id
// 	container, err := tagService.TagRepository.GetByID(tagID)
// 	return container, err
// }

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
