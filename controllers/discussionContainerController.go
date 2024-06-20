package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type DiscussionContainerController struct {
	DiscussionContainerService interfaces.DiscussionContainerService
}

// GetDiscussionContainer godoc
// @Summary 	Get discussion container
// @Description Get a discussion container by its ID, to access its discussions
// @Tags 		discussion-containers
// @Accept  	json
// @Param		discussionContainerID		path		string			true	"Discussion Container ID"
// @Produce		json
// @Success 	200 		{object}	models.DiscussionContainerDTO
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Failure		500			{object} 	utils.HTTPError
// @Router 		/discussion-containers/{discussionContainerID}	[get]
func (discussionContainerController *DiscussionContainerController) GetDiscussionContainer(c *gin.Context) {
	// Get the discussion container ID from the path
	discussionContainerIDString := c.Param("discussionContainerID")

	discussionContainerID, err := strconv.ParseUint(discussionContainerIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse discussion container ID '%s' as unsigned integer: %s", discussionContainerIDString, err)})

		return
	}

	// Fetch the discussion container from the database
	discussionContainer, err := discussionContainerController.DiscussionContainerService.GetDiscussionContainer(uint(discussionContainerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not get discussion container: %v", err.Error())})

		return
	}

	c.JSON(http.StatusOK, discussionContainer)
}
