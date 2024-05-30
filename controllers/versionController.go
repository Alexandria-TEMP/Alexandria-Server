package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v2

type VersionController struct {
	VersionService interfaces.VersionService
}

// GetVersion godoc
// @Summary 	Get version
// @Description Get a version by version ID
// @Tags 		versions
// @Accept  	json
// @Param		versionID		path		string			true	"Version ID"
// @Produce		json
// @Success 	200 		{object}	models.VersionDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/versions/{versionID}	[get]
func (versionController *VersionController) GetVersion(_ *gin.Context) {

}

// CreateVersion godoc
// @Summary 	Create new version
// @Description Create a new version with discussions and repository
// @Tags 		versions
// @Accept  	multipart/form-data
// @Param		postID		query		string					true	"Parent Post ID"
// @Param		repository	body		forms.IncomingFileForm	true	"Repository to create"
// @Produce		application/json
// @Success 	200		{object}	models.VersionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions/{postID}	[post]
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

// GetRender godoc
// @Summary Get the render of a version
// @Description Get the render of the repository underlying a version
// @Tags 		versions
// @Param		versionID		path		string			true	"version ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/render	[get]
func (versionController *VersionController) GetRender(_ *gin.Context) {
	// TODO: find out how to send back html file in godoc
}

// GetRepository godoc specs
// @Summary 	Get the repository of a version
// @Description Get the entire zipped repository of a version
// @Tags 		versions
// @Param		versionID	path		string				true	"Version ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions/{versionID}/repository	[get]
func (versionController *VersionController) GetRepository(_ *gin.Context) {

}

// GetFileTree godoc
// @Summary 	Get the file tree of a repository
// @Description Get the file tree of a repository of a version
// @Tags 		versions
// @Accept  	json
// @Param		versionID		path		string			true	"version ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/tree	[get]
func (versionController *VersionController) GetFileTree(_ *gin.Context) {

}

// GetFileFromrepository godoc
// @Summary 	Get a file from a repository
// @Description Get the contents of a single file from a repository of a version
// @Tags 		versions
// @Param		versionID		path		string			true	"version ID"
// @Param		filePath		body		string			true	"file path"
// @Success 	200		{object}	[]byte
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/file	[get]
func (versionController *VersionController) GetFileFromrepository(_ *gin.Context) {
	// TODO: find out if this response type is correct
}
