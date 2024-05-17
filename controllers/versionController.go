package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
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
// @Success 	200		{object}	models.Version
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/version/{postID}	[post]
func (versionController *VersionController) CreateVersion(c *gin.Context) {
	// extract file
	incomingFileForm := forms.IncomingFileForm{}
	err := c.ShouldBindWith(&incomingFileForm, binding.FormMultipart)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind IncomingFileForm from request body"))

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
	version, err := versionController.VersionService.CreateVersion(c, incomingFileForm.File, uint(postID))
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusInternalServerError, fmt.Errorf("%w", err))

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, version)
}
