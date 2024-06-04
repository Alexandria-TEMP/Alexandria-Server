package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type DiscussionController struct {
}

// GetDiscussion godoc
// @Summary 	Get discussion
// @Description Get a discussion by discussion ID
// @Accept  	json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 		{object}	models.DiscussionDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/discussions/{discussionID}	[get]
func (discussionController *DiscussionController) GetDiscussion(_ *gin.Context) {

}

// CreateDiscussion godoc
// @Summary 	Create new discussion
// @Description Create a new discussion
// @Description If parent ID field is used, the discussion will be a reply
// @Accept  	json
// @Param		form	body	forms.DiscussionCreationForm	true	"Discussion Creation Form"
// @Param 		parentID			query		string			false	"Parent ID"
// @Produce		json
// @Success 	200 	{object} 	models.DiscussionDTO
// @Failure		400
// @Failure		500
// @Router 		/discussions 		[post]
func (discussionController *DiscussionController) CreateDiscussion(_ *gin.Context) {

}

// DeleteDiscussion godoc
// @Summary 	Delete a discussion
// @Description Delete a discussion with given ID from database
// @Accept  	json
// @Param		discussionID		path		string			true	"discussion ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/discussions/{discussionID} 		[delete]
func (discussionController *DiscussionController) DeleteDiscussion(_ *gin.Context) {
	// delete method goes here
}

// GetDiscussionReplies godoc
// @Summary 	Get all the replies of a discussion
// @Description Gets an array of all the first-level replies of a discussion
// @Description Endpoint is offset-paginated
// @Accept  	json
// @Param		discussionID		path		string			true	"discussion ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/discussions/{discussionID}/replies 		[get]
func (discussionController *DiscussionController) GetDiscussionReplies(_ *gin.Context) {
	//TODO: make paginated
}

// AddDiscussionReport godoc
// @Summary 	Add a new report to a discussion
// @Description Create a new report for a discussion
// @Accept  	json
// @Param		form	body	forms.ReportCreationForm	true	"Report Creation Form"
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 	{object} 	models.ReportDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/discussions/{discussionID}/reports 		[post]
func (discussionController *DiscussionController) AddDiscussionReport(_ *gin.Context) {

}

// GetDiscussionReports godoc
// @Summary		Get all reports of this discussion
// @Description	Get all reports that have been added to this discussion
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.ReportDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/discussions/{discussionID}/reports 		[get]
func (discussionController *DiscussionController) GetDiscussionReports(_ *gin.Context) {
	//TODO: make paginated
}
