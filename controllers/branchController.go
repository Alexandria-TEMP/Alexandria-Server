package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type BranchController struct {
	BranchService interfaces.BranchService
}

// GetBranch godoc
// @Summary 	Get branch
// @Description Get a branch by branch ID
// @Accept  	json
// @Param		branchID		path		string			true	"Branch ID"
// @Produce		json
// @Success 	200 		{object}	models.BranchDTO
// @Failure		400 		{object}
// @Failure		404 		{object}
// @Failure		500 		{object}
// @Router 		/branches/{branchID}	[get]
func (branchController *BranchController) GetBranch(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%s ", branchIDStr)})

		return
	}

	// get branch and check it exists
	branch, err := branchController.BranchService.GetBranch(uint(branchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot find any branch with id=%s ", branchIDStr)})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, branch.IntoDTO)
}

// CreateBranch godoc
// @Summary 	Create new branch
// @Description Create a new question or discussion branch
// @Accept  	json
// @Param		form	body	forms.BranchCreationForm	true	"Branch Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.BranchDTO
// @Failure		404 	{object}
// @Failure		500 	{object}
// @Router 		/branches 		[post]
func (branchController *BranchController) CreateBranch(c *gin.Context) {
	// extract branchCreationForm
	form := forms.BranchCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind BranchCreationForm from request body"})

		return
	}

	branch, err404, err500 := branchController.BranchService.CreateBranch(form)

	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, branch.IntoDTO)
}

// GetReviewStatus godoc
// @Summary 	Returns status of all branch reviews
// @Description Returns an array of the statuses of all the reviews of this branch
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200		{array}		string
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router 		/branches/{branchID}/reviews		[get]
func (branchController *BranchController) GetReviewStatus(_ *gin.Context) {
	// delete method goes here
}

// GetReview godoc
// @Summary 	Returns a review of a branch
// @Description Returns a review with the given ID of the branch with the given ID
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		reviewID			path		string			true	"review ID"
// @Produce		json
// @Success 	200		{object}	models.ReviewDTO
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router 		/branches/{branchID}/reviews/{reviewID}		[get]
func (branchController *BranchController) GetReview(_ *gin.Context) {

}

// CreateReview godoc
// @Summary 	Adds a review to a branch
// @Description Adds a review to a branch
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		form	body	forms.ReviewCreationForm	true	"review creation form"
// @Produce		json
// @Success 	200
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router 		/branches/{branchID}/reviews		[post]
func (branchController *BranchController) CreateReview(_ *gin.Context) {

}

// UserCanReview godoc
// @Summary 	Returns whether the user is allowed to review this branch
// @Description Returns true if the user fulfills the requirements to review the branch
// @Description Returns false if user is unauthorized to review the branch
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		userID			path		string			true	"user ID"
// @Produce		json
// @Success 	200		{array}		boolean
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router 		/branches/{branchID}/can-review/{userID}		[get]
func (branchController *BranchController) UserCanReview(_ *gin.Context) {

}

// GetRender
// @Summary 	Get the render of a version
// @Description Get the render of the repository underlying a version if it exists and has been rendered successfully
// @Param		versionID	path		string				true	"Version ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Success		202
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Router 		/{versionID}/render	[get]
func (branchController *BranchController) GetRender(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	// get render filepath
	filePath, err202, err404 := branchController.VersionService.GetRenderFile(uint(versionID))

	// if render is pending return 202 accepted
	if err202 != nil {
		c.Status(http.StatusAccepted)
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
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500 	{object}
// @Router 		/{versionID}/repository	[get]
func (branchController *BranchController) GetRepository(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	// get repository filepath
	filePath, err := branchController.VersionService.GetRepositoryFile(uint(versionID))

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

// CreateVersion
// @Summary 	Create new version
// @Description Create a new version with discussions and repository from zipped file in body
// @Accept  	multipart/form-data
// @Param		fromVersionID		path		string			true	"Version ID"
// @Param		repository			body		file				true	"Repository to create"
// @Produce		application/json
// @Success 	200		{object}	models.Version
// @Failure		400 	{object}
// @Failure		500 	{object}
// @Router 		/{fromVersionID}		[post]
func (branchController *BranchController) CreateVersion(c *gin.Context) {
	// extract file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file found"})

		return
	}

	// extract version id
	fromVersionIDstr := c.Param("fromVersionIDs")
	fromVersionID, err := strconv.ParseUint(fromVersionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", fromVersionID)})

		return
	}

	// Create Version
	version, err := branchController.VersionService.CreateVersion(c, file, uint(fromVersionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create version"})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, version)
}

// GetFileTree godoc specs are subject to change
// @Summary 	Get the file tree of a repository
// @Description Get the file tree of a repository of a version
// @Param		versionID	path		string				true	"Version ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500 	{object}
// @Router 		/{versionID}/tree		[get]
func (branchController *BranchController) GetFileTree(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	fileTree, err1, err2 := branchController.VersionService.GetTreeFromRepository(uint(versionID))

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
// @Failure		404 	{object}
// @Failure		500 	{object}
// @Router 		/{versionID}/file/{filepath}	[get]
func (branchController *BranchController) GetFileFromRepository(c *gin.Context) {
	// extract version id
	versionIDstr := c.Param("versionID")
	versionID, err := strconv.ParseUint(versionIDstr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid version ID, cannot interpret as integer, id=%v ", versionIDstr)})

		return
	}

	relFilepath := c.Param("filepath")
	absFilepath, err := branchController.VersionService.GetFileFromRepository(uint(versionID), relFilepath)

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
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router		/{versionID}/discussions 	[get]
func (branchController *BranchController) GetDiscussions(_ *gin.Context) {

}
