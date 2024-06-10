package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

// run to create the mock
//go:generate mockgen -package=mocks -source=./tagService_interface.go -destination=../../mocks/tagService_mock.go

type TagService interface {
	GetAllScientificFieldTags() ([]*tags.ScientificFieldTag, error)
	GetTagsFromIDs(_ []uint) ([]*tags.ScientificFieldTag, error)
}
