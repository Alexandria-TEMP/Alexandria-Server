package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type FilterController struct {
}

// FilterPosts godoc
// @Summary 	Filters all posts
// @Description Returns all post IDs that meet the requirements in the form
// @Description Endpoint is offset-paginated
// @Tags 		filtering
// @Accept  	json
// @Param		form		body		forms.FilterForm	true	"Filter form"
// @Param 		page		query		uint				false	"page query"
// @Param		pageSize	query		uint				false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/filter/posts		[get]
func (filterController *FilterController) FilterPosts(_ *gin.Context) {

}

// FilterProjectPosts godoc
// @Summary 	Filters all project posts
// @Description Returns all project post IDs that meet the requirements in the form
// @Description Endpoint is offset-paginated
// @Tags 		filtering
// @Accept  	json
// @Param		form		body		forms.FilterForm	true	"Filter form"
// @Param 		page		query		uint				false	"page query"
// @Param		pageSize	query		uint				false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/filter/project-posts		[get]
func (filterController *FilterController) FilterProjectPosts(_ *gin.Context) {

}
