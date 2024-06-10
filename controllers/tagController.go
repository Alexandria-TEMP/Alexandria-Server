package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tags "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v2

type TagController struct {
	TagService interfaces.TagService
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
