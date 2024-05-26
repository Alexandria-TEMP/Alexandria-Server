package controllers

import "github.com/gin-gonic/gin"

type DiscussionController struct {
}

// GetDiscussion godoc
// @Summary 	Get discussion
// @Description Get a discussion by discussion ID
// @Accept  	json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 		{object}	models.DiscussionDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}	[get]
func (versionController *VersionController) GetDiscussion(c *gin.Context) {

}

// CreateDiscussion godoc
// @Summary 	Create new discussion
// @Description Create a new discussion
// @Description If parent ID field is used, the discussion will be a reply
// @Accept  	json
// @Param		form	body	forms.DiscussionCreationForm	true	"Discussion Creation Form"
// @Param 		parentID			body		string			false	"Parent ID"
// @Produce		json
// @Success 	200 	{object} 	models.DiscussionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/discussions 		[post]
func (versionController *VersionController) CreateDiscussion(c *gin.Context) {

}

// DeleteDiscussion godoc
// @Summary 	Delete a discussion
// @Description Delete a discussion with given ID from database
// @Accept  	json
// @Param		discussionID		path		string			true	"discussion ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/discussions/{discussionID} 		[delete]
func (versionController *VersionController) DeleteDiscussion(c *gin.Context) {
	//delete method goes here
}

// GetDiscussionReplies godoc
// @Summary 	Get all the replies of a discussion
// @Description Gets an array of all the first-level replies of a discussion
// @Accept  	json
// @Param		discussionID		path		string			true	"discussion ID"
// @Produce		json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/discussions/{discussionID}/replies 		[get]
func (versionController *VersionController) GetDiscussionReplies(c *gin.Context) {
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}/reports 		[post]
func (versionController *VersionController) AddDiscussionReport(c *gin.Context) {

}

// GetDiscussionReports godoc
// @Summary		Get all reports of this discussion
// @Description	Get all reports that have been added to this discussion
// @Accept 		json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200		{array}		models.ReportDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/discussions/{discussionID}/reports 		[get]
func (versionController *VersionController) GetDiscussionReports(c *gin.Context) {
	//TODO: make paginated
}

// - `/discussions`
//   - `POST` (mandatory version ID, optional parent discussion ID)
//   - `/:id` `GET` (for nested discussions)
//   - `/:id/replies` `GET` all discussions with that parent id
//   - `/:id` `DELETE`
//   - `/:id/reports` `POST` (wahhh)
//   - `/:id/reports` `GET` **_p_**