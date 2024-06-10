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
// @Param		form	body		forms.FilterForm	true	"Filter form"
// @Param 		page	query		uint				false	"page query"
// @Param		size	query		uint				false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/filter/posts		[get]
func (filterController *FilterController) FilterPosts(c *gin.Context) {
	var filterForm forms.FilterForm

	err := c.BindJSON(&filterForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind form JSON: %s", err)})

		return
	}

	if !filterForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	page := c.GetInt("page")
	size := c.GetInt("size")

	postIDs, err := filterController.PostService.Filter(page, size, filterForm)
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
// @Param		form		body		forms.FilterForm	true	"Filter form"
// @Param 		page		query		uint				false	"page query"
// @Param		size		query		uint				false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/filter/project-posts		[get]
func (filterController *FilterController) FilterProjectPosts(c *gin.Context) {
	var filterForm forms.FilterForm

	err := c.BindJSON(&filterForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind form JSON: %s", err)})

		return
	}

	if !filterForm.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	page := c.GetInt("page")
	size := c.GetInt("size")

	projectPostIDs, err := filterController.ProjectPostService.Filter(page, size, filterForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("filtering project posts failed: %s", err)})
	}

	c.JSON(http.StatusOK, projectPostIDs)
}
