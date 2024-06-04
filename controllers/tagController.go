package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type TagController struct {
}

// FilterPosts godoc
// @Summary 	Returns all scientific tags
// @Description Returns all scientific tags in the database
// @Produce		json
// @Success 	200		{array}		tags.ScientificFieldTag
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/tags/scientific	[get]
func (tagController *TagController) GetScientificTags(_ *gin.Context) {

}
