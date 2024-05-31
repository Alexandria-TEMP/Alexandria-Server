package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type DiscussionController struct {
}

// GetDiscussion godoc
// @Summary 	Get discussion
// @Description Get a discussion by discussion ID
// @Tags 		discussions
// @Accept  	json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 		{object}	models.DiscussionDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}	[get]
func (discussionController *DiscussionController) GetDiscussion(_ *gin.Context) {

}

// CreateDiscussion godoc
// @Summary 	Create new discussion
// @Description Create a new discussion
// @Description Either parent ID or version ID must be specified. This determines whether it's a reply or not, respectively.
// @Tags 		discussions
// @Accept  	json
// @Param		form	body	forms.DiscussionCreationForm	true	"Discussion Creation Form"
// @Param 		parentID			query		string			false	"Parent ID"
// @Param 		versionID			query		string			false	"Version ID"
// @Produce		json
// @Success 	200 	{object} 	models.DiscussionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/discussions 		[post]
func (discussionController *DiscussionController) CreateDiscussion(_ *gin.Context) {

}

// DeleteDiscussion godoc
// @Summary 	Delete a discussion
// @Description Delete a discussion with given ID from database
// @Tags 		discussions
// @Accept  	json
// @Param		discussionID		path		string			true	"discussion ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/discussions/{discussionID} 		[delete]
func (discussionController *DiscussionController) DeleteDiscussion(_ *gin.Context) {
	// delete method goes here
}

// AddDiscussionReport godoc
// @Summary 	Add a new report to a discussion
// @Description Create a new report for a discussion
// @Tags 		discussions
// @Accept  	json
// @Param		form	body	forms.ReportCreationForm	true	"Report Creation Form"
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 	{object} 	models.ReportDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}/reports 		[post]
func (discussionController *DiscussionController) AddDiscussionReport(_ *gin.Context) {

}

// GetDiscussionReports godoc
// @Summary		Get all reports of this discussion
// @Description	Get all reports that have been added to this discussion
// @Description Endpoint is offset-paginated
// @Tags 		discussions
// @Accept 		json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/discussions/{discussionID}/reports 		[get]
func (discussionController *DiscussionController) GetDiscussionReports(_ *gin.Context) {
	//TODO: make paginated
}

// GetDiscussionReport godoc
// @Summary		Gets a discussion report by ID
// @Description	Gets a discussion report by its ID
// @Tags		discussions
// @Param		reportID	path	string	true	"Report ID"
// @Produce		json
// @Success		200		{object}	reports.DiscussionReportDTO
// @Failure		400		{object}	utils.HTTPError
// @Failure		404		{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/discussions/reports/{reportID}				[get]
func (discussionController *DiscussionController) GetDiscussionReport(_ *gin.Context) {
	// TODO implement
}
