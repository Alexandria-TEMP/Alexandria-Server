package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type TagController struct {
	TagService interfaces.TagService
}

// GetScientificFieldTag godoc
// @Summary 	Get scientific field tag from database
// @Description Get a scientific field tag by tag ID
// @Tags 		tags
// @Accept  	json
// @Param		tagID		path		string			true	"tag ID"
// @Produce		json
// @Success 	200 	{object}	models.ScientificFieldTagDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/tags/scientific/:tagID	[get]
func (tagController *TagController) GetScientificFieldTag(c *gin.Context) {
	// extract the id of the scientific field tag
	tagIDStr := c.Param("tagID")
	initTagID, err := strconv.ParseUint(tagIDStr, 10, 64)

	// if this caused an error, print it and return status 400: bad input
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid tag ID, cannot interpret '%s' as integer: %s ", tagIDStr, err)})

		return
	}

	// cast tag ID as uint instead of uint64, because database only accepts those
	tagID := uint(initTagID)

	// get the tag through the service
	tag, err := tagController.TagService.GetTagByID(tagID)

	// if there was an error, print it and return status 404: not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot get member because no tag with ID '%d' exists: %s", tagID, err)})

		return
	}

	// if correct response send the tag back
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, tag)
}

// GetScientificTags godoc
// @Summary 	Returns all scientific tags
// @Description Returns all scientific tags in the database
// @Tags 		tags
// @Produce		json
// @Success 	200		{array}		models.ScientificFieldTagDTO
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/tags/scientific	[get]
func (tagController *TagController) GetScientificTags(c *gin.Context) {
	tagObjects, err := tagController.TagService.GetAllScientificFieldTags()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot get tags: %s", err)})

		return
	}

	tagDTOs := []models.ScientificFieldTagDTO{}

	for _, tag := range tagObjects {
		dto := tag.IntoDTO()

		tagDTOs = append(tagDTOs, dto)
	}

	// if correct response send the tags back
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, tagDTOs)
}

// GetCompletionStatusTags godoc
// @Summary 	Returns all completion statuses
// @Description Returns every possible completion status that a Post can have
// @Tags		tags
// @Produce		json
// @Success		200		{array}		models.ProjectCompletionStatus
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/tags/completion-status	[get]
func (tagController *TagController) GetCompletionStatusTags(c *gin.Context) {
	completionStatusTags := []models.PostType{models.Project, models.Question, models.Reflection}

	c.JSON(http.StatusOK, completionStatusTags)
}

// GetPostTypeTags godoc
// @Summary 	Returns all post types
// @Description Returns every possible post type that a Post can have
// @Tags		tags
// @Produce		json
// @Success		200		{array}		models.PostType
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/tags/post-type	[get]
func (tagController *TagController) GetPostTypeTags(c *gin.Context) {
	postTypeTags := []models.PostType{models.Project, models.Question, models.Reflection}

	c.JSON(http.StatusOK, postTypeTags)
}

// GetFeedbackPreferenceTags godoc
// @Summary 	Returns all feedback preferences
// @Description Returns every possible feedback preference that a Project Post can have
// @Tags		tags
// @Produce		json
// @Success		200		{array}		models.ProjectFeedbackPreference
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/tags/feedback-preference	[get]
func (tagController *TagController) GetFeedbackPreferenceTags(c *gin.Context) {
	feedbackPreferenceTags := []models.ProjectFeedbackPreference{models.DiscussionFeedback, models.FormalFeedback}

	c.JSON(http.StatusOK, feedbackPreferenceTags)
}
