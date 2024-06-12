package controllers

import (
	"fmt"
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
// @Param		form	body		forms.PostFilterForm	true	"Post filter form"
// @Param 		page	query		uint					false	"page query"
// @Param		size	query		uint					false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/filter/posts		[get]
func (filterController *FilterController) FilterPosts(c *gin.Context) {
	var postFilterForm forms.PostFilterForm

	err := c.BindJSON(&postFilterForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind form JSON: %s", err)})

		return
	}

	if !postFilterForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	page := c.GetInt("page")
	size := c.GetInt("size")

	postIDs, err := filterController.PostService.Filter(page, size, postFilterForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("filtering posts failed: %s", err)})

		return
	}

	c.JSON(http.StatusOK, postIDs)
}

// FilterProjectPosts godoc
// @Summary 	Filters all project posts
// @Description Returns all project post IDs that meet the requirements in the form
// @Description Endpoint is offset-paginated
// @Tags 		filtering
// @Accept  	json
// @Param		form		body		forms.ProjectPostFilterForm	true	"Project post filter form"
// @Param 		page		query		uint						false	"page query"
// @Param		size		query		uint						false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/filter/project-posts		[get]
func (filterController *FilterController) FilterProjectPosts(c *gin.Context) {
	var projectPostFilterForm forms.ProjectPostFilterForm

	err := c.BindJSON(&projectPostFilterForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind form JSON: %s", err)})

		return
	}

	if !projectPostFilterForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	page := c.GetInt("page")
	size := c.GetInt("size")

	projectPostIDs, err := filterController.ProjectPostService.Filter(page, size, projectPostFilterForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("filtering project posts failed: %s", err)})
	}

	c.JSON(http.StatusOK, projectPostIDs)
}
