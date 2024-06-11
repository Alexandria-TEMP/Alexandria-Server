package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	tags "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v2

type TagController struct {
	TagService interfaces.TagService
}

// GetScientificFieldTag godoc
// @Summary 	Get scientific field tag from database
// @Description Get a scientific field tag by tag ID
// @Tags 		members
// @Accept  	json
// @Param		tagID		path		string			true	"tag ID"
// @Produce		json
// @Success 	200 		{object}	models.ScientificFieldTagDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500			{object}	utils.HTTPError
// @Router 		/members/{memberID}	[get]
// @Router 		/tags/scientific/:tagID	[get]
func (tagController *TagController) GetScientificFieldTag(c *gin.Context) {
	// extract the id of the scientific field tag
	tagIDStr := c.Param("tagID")
	initTagID, err := strconv.ParseUint(tagIDStr, 10, 64)

	// if this caused an error, print it and return status 400: bad input
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid tag ID, cannot interpret as integer, id=%s ", tagIDStr))

		return
	}

	// cast tag ID as uint instead of uint64, because database only accepts those
	tagID := uint(initTagID)

	// get the tag through the service
	tag, err := tagController.TagService.GetTagByID(tagID)

	// if there was an error, print it and return status 404: not found
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusNotFound, fmt.Errorf("cannot get member because no tag with this ID exists, id=%d", tagID))

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
// @Success 	200		{array}		tags.ScientificFieldTagDTO
// @Failure		404 	{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/tags/scientific	[get]
func (tagController *TagController) GetScientificTags(c *gin.Context) {
	tagObjects, err := tagController.TagService.GetAllScientificFieldTags()

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusNotFound, fmt.Errorf("cannot get tags, error: %w", err))

		return
	}

	tagDTOs := []tags.ScientificFieldTagDTO{}

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
// @Failure		400 	{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/tags/completion-status	[get]
func (tagController *TagController) GetCompletionStatusTags(_ *gin.Context) {
	// TODO implement
}

// GetPostTypeTags godoc
// @Summary 	Returns all post types
// @Description Returns every possible post type that a Post can have
// @Tags		tags
// @Produce		json
// @Success		200		{array}		models.PostType
// @Failure		400 	{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/tags/post-type	[get]
func (tagController *TagController) GetPostTypeTags(_ *gin.Context) {
	// TODO implement
}

// GetFeedbackPreferenceTags godoc
// @Summary 	Returns all feedback preferences
// @Description Returns every possible feedback preference that a Project Post can have
// @Tags		tags
// @Produce		json
// @Success		200		{array}		models.ProjectFeedbackPreference
// @Failure		400 	{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/tags/feedback-preference	[get]
func (tagController *TagController) GetFeedbackPreferenceTags(_ *gin.Context) {
	// TODO implement
}
