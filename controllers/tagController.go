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
// @Success		200		{array}			tags.CompletionStatus
// @Failure		500		{object}		utils.HTTPError
// @Router		/tags/completion-status	[get]
func (tagController *TagController) GetCompletionStatusTags(_ *gin.Context) {
	// TODO implement
}
