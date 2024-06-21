package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type FilterController struct {
	PostService interfaces.PostService
}

// FilterPosts godoc
// @Summary 	Filters all posts
// @Description Returns all post IDs that meet the requirements in the form
// @Description Endpoint is offset-paginated
// @Tags 		filtering
// @Accept  	json
// @Param 		page	query		uint					false	"page query"
// @Param		size	query		uint					false	"page size"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/filter/posts		[get]
func (filterController *FilterController) FilterPosts(c *gin.Context) {
	page := c.GetInt("page")
	size := c.GetInt("size")

	postIDs, err := filterController.PostService.Filter(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("filtering posts failed: %v", err.Error())})

		return
	}

	c.JSON(http.StatusOK, postIDs)
}
