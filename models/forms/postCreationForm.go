package forms

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type PostCreationForm struct {
	models.PostMetadata
	CurrentVersion models.Version
}
