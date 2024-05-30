package forms

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"

type PostCreationForm struct {
	// TODO send files somehow?

	AuthorMemberIDs     []uint // Members that are authors of the post
	Title               string
	PostType            tags.PostType
	ScientificFieldTags []tags.ScientificField
}
