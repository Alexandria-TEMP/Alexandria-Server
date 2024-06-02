package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

const headerSize = 512

type VersionController struct {
	VersionService interfaces.VersionService
}

// GetVersion
// @Summary 	Get version
// @Description Get a version by version ID
// @Param		versionID		path		string			true	"Version ID"
// @Produce		application/json
// @Success 	200 		{object}	models.VersionDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Router 		/versions/{versionID}	[get]
func (versionController *VersionController) GetVersion(_ *gin.Context) {

}

// CreateVersion
// @Summary 	Create new version
// @Description Create a new version with discussions and repository from zipped file in body
// @Accept  	multipart/form-data
// @Param		repository			formData		file				true	"Repository to create"
// @Produce		application/json
// @Success 	200		{object}	models.VersionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions		[post]
func (versionController *VersionController) CreateVersion(c *gin.Context) {
	// extract file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file found"})

		return
	}

	// Create Version
	version, err := versionController.VersionService.CreateVersion(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create version"})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, version.IntoDTO())
}

// GetRender
// @Summary 	Get the render of a version
// @Description Get the render of the repository underlying a version if it exists and has been rendered successfully
// @Param		versionID	path		string				true	"Version ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Success		202		{object}	string
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Router 		/versions/{versionID}/render	[get]
func (versionController *VersionController) GetRender(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	// get render filepath
	filePath, err202, err404 := versionController.VersionService.GetRenderFile(uint(versionID))

	// if render is pending return 202 accepted
	if err202 != nil {
		c.String(http.StatusAccepted, "text/plain", "pending")
		return
	}

	// if render is failed return 404 not found
	if err404 != nil {
		c.JSON(http.StatusNotFound, err404)

		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=render.html")
	c.Header("Content-Type", "text/html")
	c.File(filePath)
}

// GetRepository godoc specs are subject to change
// @Summary 	Get the repository of a version
// @Description Get the entire zipped repository of a version
// @Param		versionID	path		string				true	"Version ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions/{versionID}/repository	[get]
func (versionController *VersionController) GetRepository(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	// get repository filepath
	filePath, err := versionController.VersionService.GetRepositoryFile(uint(versionID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no such repository found"})

		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=quarto_project.zip")
	c.Header("Content-Type", "application/zip")
	c.File(filePath)
}

// GetFileTree godoc specs are subject to change
// @Summary 	Get the file tree of a repository
// @Description Get the file tree of a repository of a version
// @Param		versionID	path		string				true	"Version ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions/{versionID}/tree		[get]
func (versionController *VersionController) GetFileTree(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	fileTree, err1, err2 := versionController.VersionService.GetTreeFromRepository(uint(versionID))

	// if repository doesnt exist throw 404 not found
	if err1 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no such repository found"})

		return
	}

	// if failed to parse file tree throw 500 internal server error
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse file tree"})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, fileTree)
}

// GetFileFromRepository godoc specs are subject to change
// @Summary 	Get a file from a repository
// @Description Get the contents of a single file from a repository of a version
// @Param		versionID	path		string				true	"Version ID"
// @Param		filepath	path		string				true	"Filepath"
// @Produce		application/octet-stream
// @Success 	200		{object}	[]byte
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions/{versionID}/file/{filepath}	[get]
func (versionController *VersionController) GetFileFromRepository(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	relFilepath := c.Param("filepath")
	absFilepath, err := versionController.VersionService.GetFileFromRepository(uint(versionID), relFilepath)

	// if files doesnt exist return 404 not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no such file exists"})

		return
	}

	fileContentType, err1 := mimetype.DetectFile(absFilepath)

	fileData, err2 := os.Open(absFilepath)

	// Get the file info
	fileInfo, err3 := fileData.Stat()

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})

		return
	}

	defer fileData.Close()

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	c.Header("Content-Type", fileContentType.String())
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(absFilepath)
}

// GetDiscussions godoc
// @Summary Returns all level 1 discussions associated with the version
// @Description Returns all discussions on this version that are not a reply to another discussion
// @Description Endpoint is offset-paginated
// @Param		versionID	path		string			true	"version ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		application/json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/discussions 	[get]
func (versionController *VersionController) GetDiscussions(_ *gin.Context) {

}
