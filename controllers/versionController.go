package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v1

type VersionController struct {
	VersionService interfaces.VersionService
}

// CreateVersion godoc
// @Summary 	Create new version
// @Description Create a new version with discussions and repository
// @Accept  	multipart/form-data
// @Param		postID		path		string				true	"Parent Post ID"
// @Param		repository	body		models.Repository	true	"Repository to create"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Router 		/version/{postID}	[post]
func (versionController *VersionController) CreateVersion(c *gin.Context) {
	// other way of doing this
	// form, err := c.MultipartForm()
	// files := form.File["file"]

	// extract file
	repository := models.Repository{}
	err := c.ShouldBindWith(&repository, binding.FormMultipart)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind Repository from request body"))

		return
	}

	// extract post id
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid article ID, cannot interpret as integer, id=%v ", postIDStr))

		return
	}

	// Create Version here
	version := versionController.VersionService.CreateVersion()

	// Create and add post to database here. For now just do this to test.
	err = versionController.VersionService.SaveRepository(c, repository.File, uint(version.ID), uint(postID))

	// response
	c.Status(http.StatusOK)
}
