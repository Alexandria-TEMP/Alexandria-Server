package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type TagController struct {
}

// FilterPosts godoc
// @Summary 	Returns all scientific tags
// @Description Returns all scientific tags (an array of strings) in the database
// @Tags 		scientific-field-tags
// @Produce		json
// @Success 	200		{array}		tags.ScientificFieldTag
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/tags/scientific	[get]
func (tagController *TagController) GetScientificTags(_ *gin.Context) {

}
