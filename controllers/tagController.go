package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type TagController struct {
}

// GetScientificTags godoc
// @Summary 	Returns all scientific tags
// @Description Returns all scientific tags (an array of strings) in the database
// @Tags 		tags
// @Produce		json
// @Success 	200		{array}		tags.ScientificFieldTag
// @Failure		400 	{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/tags/scientific	[get]
func (tagController *TagController) GetScientificTags(_ *gin.Context) {
	// TODO implement
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
