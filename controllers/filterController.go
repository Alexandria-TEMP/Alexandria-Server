package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type FilterController struct {
	PostService        interfaces.PostService
	ProjectPostService interfaces.ProjectPostService
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
func (filterController *FilterController) FilterPosts(c *gin.Context) {
	var filterForm forms.FilterForm
	c.BindJSON(&filterForm)

	if !filterForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	// TODO
	// postIDs := filterController.PostService.filter(filterForm)

	// c.JSON(http.StatusOK, postIDs)
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
