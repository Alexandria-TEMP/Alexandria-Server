package controllers

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v2

type VersionController struct {
}

// GetVersion godoc
// @Summary 	Get version
// @Description Get a version by version ID
// @Accept  	json
// @Param		versionID		path		string			true	"Version ID"
// @Produce		json
// @Success 	200 		{object}	models.VersionDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/versions/{versionID}	[get]
func (versionController *VersionController) GetVersion(c *gin.Context) {

}


// CreateVersion godoc
// @Summary 	Create new version
// @Description Create a new version with discussions and repository
// @Accept  	multipart/form-data
// @Param		postID		query		string				true	"Parent Post ID"
// @Param		repository	body		models.Repository	true	"Repository to create"
// @Produce		application/json
// @Success 	200		{object}	models.VersionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions	[post]
func (versionController *VersionController) CreateVersion(c *gin.Context) {

}


// RenderVersion godoc
// @Summary Get the render of a version
// @Description Get the render of the repository underlying a version
// @Param		versionID		path		string			true	"version ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/render	[get]
func (versionController *VersionController) RenderVersion(c *gin.Context) {
	//TODO: find out how to send back html file in godoc
}


// GetRepository godoc specs
// @Summary 	Get the repository of a version
// @Description Get the entire zipped repository of a version
// @Param		versionID	path		string				true	"Version ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/version/{versionID}/repository	[get]
func (versionController *VersionController) GetRepository(c *gin.Context) {

}

// GetFileTreeVersion godoc
// @Summary 	Get the file tree of a repository
// @Description Get the file tree of a repository of a version
// @Accept  	json
// @Param		versionID		path		string			true	"version ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/tree	[get]
func (versionController *VersionController) GetFileTreeVersion(c *gin.Context) {

}



// GetFileFromVersion godoc
// @Summary 	Get a file from a repository
// @Description Get the contents of a single file from a repository of a version
// @Param		versionID		path		string			true	"version ID"
// @Param		filePath		body		string			true	"file path"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/file	[get]
func (versionController *VersionController) GetFileFromVersion(c *gin.Context) {
	//TODO: find out if this response type is correct
}



// GetVersionDiscussions godoc
// @Summary Returns all level 1 discussions associated with the version
// @Description Returns all discussions on this version that are not a reply to another discussion
// @Description Endpoint is offset-paginated
// @Accept  	json
// @Param		versionID		path		string			true	"version ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/discussions 	[get]
func (versionController *VersionController) GetVersionDiscussions(c *gin.Context) {

}





// - `/versions`
//   - `POST`
//   - `/:id` `GET` (get the files?)
//   - `/:id/file` (get one specific file)
//   - `/:id/tree`
//   - `/:id` `/render` `GET`
//   - `/:id/discussions` `GET` (gets all level-1 discussions of the version) **_p_**


