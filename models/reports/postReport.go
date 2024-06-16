package reports

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gorm.io/gorm"
)

type PostReport struct {
	gorm.Model

	// TODO implement reports

	// PostReport belongs to Post
	Post   models.Post
	PostID uint
}

type PostReportDTO struct {
	PostID uint `json:"postID" example:"1"`
}

func (model *PostReport) GetID() uint {
	return model.ID
}

func (model *PostReport) IntoDTO() PostReportDTO {
	return PostReportDTO{
		PostID: model.PostID,
	}
}

func (model *PostReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
