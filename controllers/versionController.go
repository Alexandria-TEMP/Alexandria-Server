package controllers

import (
	"github.com/gin-gonic/gin"
)

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
// @Description Create a new version
// @Accept  	json
// @Param		form	body	forms.VersionCreationForm	true	"Version Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.VersionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/versions 		[post]
func (versionController *VersionController) CreateVersion(c *gin.Context) {

}

// RenderVersion godoc
// @Summary Returns the rendered form of the version files
// @Description Returns html file of the rendered version
// @Accept  	json
// @Param		versionID		path		string			true	"version ID"
// @Produce		html
// @Success 	200		{object}	file
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/render	[get]
func (versionController *VersionController) RenderVersion(c *gin.Context) {
	//TODO: find out how to send back html file in godoc
}

// GetFileTreeVersion godoc
// @Summary Returns the base level names the version files
// @Description Returns the top layer of names of the version files
// @Accept  	json
// @Param		versionID		path		string			true	"version ID"
// @Produce		application/zip
// @Success 	200		{object}	file
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/versions/{versionID}/tree	[get]
func (versionController *VersionController) GetFileTreeVersion(c *gin.Context) {
	//TODO: find out how to send back html file in godoc
}

// GetFileFromVersion godoc
// @Summary Returns a specific file from the version repository
// @Description Returns the file from the given path
// @Accept  	json
// @Param		versionID		path		string			true	"version ID"
// @Param		filePath		body		string			true	"file path"
// @Produce		application/zip
// @Success 	200		{object}	file
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