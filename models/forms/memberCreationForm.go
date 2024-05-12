package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type MemberCreationForm struct {
	models.FirstName
	models.LastName
	models.Email
	models.Password
	models.Institution
	Posts models.[]Post
	Discussions models.[]Discussion
	Reviews models.[]MergeRequestReview
}
