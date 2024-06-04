package reports

import (
	"encoding/json"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gorm.io/gorm"
)

type DiscussionReport struct {
	gorm.Model

	// TODO implement reports

	// DiscussionReport belongs to Discussion
	Discussion   models.Discussion
	DiscussionID uint
}

type DiscussionReportDTO struct {
	DiscussionID uint
}

func (model *DiscussionReport) GetID() uint {
	return model.ID
}

func (model *DiscussionReport) IntoDTO() DiscussionReportDTO {
	return DiscussionReportDTO{
		DiscussionID: model.DiscussionID,
	}
}

func (model *DiscussionReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.IntoDTO())
}
