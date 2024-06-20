package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type DiscussionController struct {
	DiscussionService interfaces.DiscussionService
}

// GetDiscussion godoc
// @Summary 	Get discussion
// @Description Get a discussion by discussion ID
// @Tags 		discussions
// @Accept  	json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 		{object}	models.DiscussionDTO
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Failure		500			{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}	[get]
func (discussionController *DiscussionController) GetDiscussion(c *gin.Context) {
	// Parse discussion ID path parameter
	discussionIDString := c.Param("discussionID")

	discussionID, err := strconv.ParseUint(discussionIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse discussion ID '%s' as unsigned integer: %s", discussionIDString, err)})

		return
	}

	// Get from database
	discussion, err := discussionController.DiscussionService.GetDiscussion(uint(discussionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get discussion with ID %d: %s", discussionID, err)})

		return
	}

	c.JSON(http.StatusOK, discussion)
}

// CreateRootDiscussion godoc
// @Summary 	Create new root discussion
// @Description Create a new root-level discussion, meaning a discussion that is not a reply.
// @Tags 		discussions
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.RootDiscussionCreationForm	true	"Root Discussion Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.DiscussionDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/discussions/roots 		[post]
func (discussionController *DiscussionController) CreateRootDiscussion(c *gin.Context) {
	// Bind discussion creation form from request
	var discussionCreationForm forms.RootDiscussionCreationForm

	if err := c.BindJSON(&discussionCreationForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not bind form from JSON: %v", err.Error())})

		return
	}

	if !discussionCreationForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	// get member
	member, exists := c.Get("currentMember")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get logged in user"})

		return
	}

	// Create discussion in the database
	createdDiscussion, err := discussionController.DiscussionService.CreateRootDiscussion(&discussionCreationForm, member.(*models.Member))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create root discussion: %v", err.Error())})

		return
	}

	c.JSON(http.StatusOK, createdDiscussion)
}

// CreateReplyDiscussion godoc
// @Summary 	Create new reply discussion
// @Description Create a new reply-type discussion, so a discussion that is a child of another discussion.
// @Tags 		discussions
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.ReplyDiscussionCreationForm	true	"Reply Discussion Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.DiscussionDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/discussions/replies 		[post]
func (discussionController *DiscussionController) CreateReplyDiscussion(c *gin.Context) {
	// Bind discussion creation form from request
	var discussionCreationForm forms.ReplyDiscussionCreationForm

	if err := c.BindJSON(&discussionCreationForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not bind form from JSON: %v", err.Error())})

		return
	}

	if !discussionCreationForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	// get member
	member, exists := c.Get("currentMember")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get logged in user"})

		return
	}

	// Create discussion in the database
	createdDiscussion, err := discussionController.DiscussionService.CreateReply(&discussionCreationForm, member.(*models.Member))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create reply discussion: %v", err.Error())})

		return
	}

	c.JSON(http.StatusOK, createdDiscussion)
}

// DeleteDiscussion godoc
// @Summary 	Delete a discussion
// @Description Delete a discussion with given ID from database
// @Tags 		discussions
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		discussionID		path		string			true	"discussion ID"
// @Produce		json
// @Success 	200
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID} 		[delete]
func (discussionController *DiscussionController) DeleteDiscussion(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// AddDiscussionReport godoc
// @Summary 	Add a new report to a discussion
// @Description Create a new report for a discussion
// @Tags 		discussions
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.ReportCreationForm	true	"Report Creation Form"
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200 	{object} 	models.ReportDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}/reports 		[post]
func (discussionController *DiscussionController) AddDiscussionReport(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetDiscussionReports godoc
// @Summary		Get all reports of this discussion
// @Description	Get all reports that have been added to this discussion
// @Tags 		discussions
// @Accept 		json
// @Param		discussionID		path		string			true	"Discussion ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/discussions/{discussionID}/reports 		[get]
func (discussionController *DiscussionController) GetDiscussionReports(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetDiscussionReport godoc
// @Summary		Gets a discussion report by ID
// @Description	Gets a discussion report by its ID
// @Tags		discussions
// @Param		reportID	path	string	true	"Report ID"
// @Produce		json
// @Success		200		{object}	reports.DiscussionReportDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/discussions/reports/{reportID}				[get]
func (discussionController *DiscussionController) GetDiscussionReport(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
