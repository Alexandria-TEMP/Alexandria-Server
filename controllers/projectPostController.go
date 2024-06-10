package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type ProjectPostController struct {
	PostService        interfaces.PostService
	ProjectPostService interfaces.ProjectPostService
	RenderService      interfaces.RenderService
}

// GetProjectPost godoc
// @Summary 	Get project post
// @Description Get a project post by ID
// @Tags 		project-posts
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.ProjectPostDTO
// @Failure		400
// @Failure		404
// @Router 		/project-posts/{postID}	[get]
func (projectPostController *ProjectPostController) GetProjectPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not interpret ID %s as unsigned integer, reason: %s", postIDStr, err)})

		return
	}

	projectPost, err := projectPostController.ProjectPostService.GetProjectPost(uint(postID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not get project post, reason: %s", err)})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, projectPost)
}

// CreateProjectPost godoc
// @Summary 	Create new project post
// @Description Create a new project post with a single open branch. Upload to this branch in order to have your post reviewed.
// @Tags 		project-posts
// @Accept  	json
// @Param		form	body		forms.ProjectPostCreationForm	true	"Project Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400
// @Failure		500
// @Router 		/project-posts		[post]
func (projectPostController *ProjectPostController) CreateProjectPost(c *gin.Context) {
	form := forms.ProjectPostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid project post creation form: %s", err)})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	projectPost, err404, err500 := projectPostController.ProjectPostService.CreateProjectPost(&form)

	if err404 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("not found: %s", err)})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("internal server error: %s", err)})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, projectPost)
}

// UpdateProjectPost godoc
// @Summary 	Update project post
// @Description Update any number of the aspects of a project post
// @Tags 		project-posts
// @Accept  	json
// @Param		post	body		models.ProjectPostDTO		true	"Updated Project Post"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Router 		/project-posts 		[put]
func (projectPostController *ProjectPostController) UpdateProjectPost(c *gin.Context) {
	// extract post
	updatedProjectPost := models.ProjectPost{}
	err := c.BindJSON(&updatedProjectPost)

	// TODO convert from project post DTO to updated project post

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind updated ProjectPost from request body"})

		return
	}

	// Update and add post to database here. For now just do this to test.
	err = projectPostController.ProjectPostService.UpdateProjectPost(&updatedProjectPost)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot update post because no ProjectPost with this ID exists"})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}

// DeleteProjectPost godoc
// @Summary 	Delete a project post
// @Description Delete a project post with given ID from database
// @Tags 		project-posts
// @Accept  	json
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/project-posts/{postID} 		[delete]
func (projectPostController *ProjectPostController) DeleteProjectPost(_ *gin.Context) {
	// delete method goes here
}

// CreateProjectPostFromGithub godoc
// @Summary 	Create new project post with the version imported from github
// @Description Create a new project post
// @Description Creates a project post in the same way as CreateProjectPost
// @Description However, the post files are imported from the given Github repository
// @Tags 		project-posts
// @Accept  	json
// @Param		form	body	forms.ProjectPostCreationForm	true	"Post Creation Form"
// @Param		url		query	string							true	"Github repository url"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400
// @Failure		500
// @Failure 	502
// @Router 		/project-posts/from-github 		[post]
func (projectPostController *ProjectPostController) CreateProjectPostFromGithub(_ *gin.Context) {

}

// GetProjectPostDiscussions godoc
// @Summary Returns all discussion IDs associated with the project post
// @Description Returns all discussion IDs on this project post over all its previous versions, instead of only the current version
// @Tags 		project-posts
// @Accept  	json
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/project-posts/{postID}/all-discussions 	[get]
func (projectPostController *ProjectPostController) GetProjectPostDiscussions(_ *gin.Context) {
	// TODO implement
}

// GetProjectPostBranchesByStatus godoc
// @Summary 	Returns branch IDs grouped by each branch status
// @Description Returns all branch IDs of this project post, grouped by each branch's branchreview status
// @Tags		project-posts
// @Accept		json
// @Param		postID	path	string	true	"post ID"
// @Produce		json
// @Success		200		{object}	forms.GroupedBranchForm
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/project-posts/{postID}/branches-by-status	[get]
func (projectPostController *ProjectPostController) GetProjectPostBranchesByStatus(_ *gin.Context) {
	// TODO implement
}

// GetMainRender
// @Summary 	Get the main render of a project post
// @Description Get the main render of the repository underlying a project post if it exists and has been rendered successfully
// @Tags 		project-posts
// @Param		projectPostID		path		string				true	"Project Post ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Success		202		{object}	[]byte
// @Failure		400
// @Failure		404
// @Router 		/project-posts/{projectPostID}/render	[get]
func (projectPostController *ProjectPostController) GetMainRender(c *gin.Context) {
	// extract projectPostID id
	projectPostIDStr := c.Param("projectPostID")
	projectPostID, err := strconv.ParseUint(projectPostIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid project post ID, cannot interpret as integer, id=%v ", projectPostIDStr)})

		return
	}

	// get project post for the post id
	projectPost, err := projectPostController.ProjectPostService.GetProjectPost(uint(projectPostID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no such project post found with id %v", projectPostID)})

		return
	}

	// get render filepath
	filePath, err202, err404 := projectPostController.RenderService.GetMainRenderFile(uint(projectPost.PostID))

	// if render is pending return 202 accepted
	if err202 != nil {
		c.String(http.StatusAccepted, "text/plain", []byte("pending"))

		return
	}

	// if render is failed return 404 not found
	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err404.Error()})

		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=render.html")
	c.Header("Content-Type", "text/html")
	c.File(filePath)
}

// GetMainProject godoc specs are subject to change
// @Summary 	Get the main repository of a project post
// @Description Get the entire zipped main repository of a project post
// @Tags 		project-posts
// @Param		projectPostID	path		string				true	"Project Post ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		400
// @Failure		404
// @Router 		/project-posts/{projectPostID}/repository	[get]
func (projectPostController *ProjectPostController) GetMainProject(c *gin.Context) {
	// extract project post id
	projectPostIDStr := c.Param("projectPostID")
	projectPostID, err := strconv.ParseUint(projectPostIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid project post ID, cannot interpret as integer, id=%v ", projectPostIDStr)})

		return
	}

	// get project post for the post id
	projectPost, err := projectPostController.ProjectPostService.GetProjectPost(uint(projectPostID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no such project post found with id %v", projectPostID)})

		return
	}

	// get repository filepath
	filePath, err := projectPostController.PostService.GetMainProject(uint(projectPost.PostID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=quarto_project.zip")
	c.Header("Content-Type", "application/zip")
	c.File(filePath)
}

// GetMainFiletree godoc specs are subject to change
// @Summary 	Get the filetree of a project post
// @Description Get the filetree of a the main version of a project post
// @Tags 		project-posts
// @Param		projectPostID	path		string				true	"Project Post ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/project-posts/{projectPostID}/tree		[get]
func (projectPostController *ProjectPostController) GetMainFiletree(c *gin.Context) {
	// extract projectPostID
	projectPostIDStr := c.Param("projectPostID")
	projectPostID, err := strconv.ParseUint(projectPostIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid project post ID, cannot interpret as integer, id=%v ", projectPostIDStr)})

		return
	}

	// get project post for the post id
	projectPost, err := projectPostController.ProjectPostService.GetProjectPost(uint(projectPostID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no such project post found with id %v", projectPostID)})

		return
	}

	fileTree, err404, err500 := projectPostController.PostService.GetMainFiletree(uint(projectPost.PostID))

	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err404.Error()})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err500.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, fileTree)
}

// GetMainFileFromProject godoc specs are subject to change
// @Summary 	Get a file from a project post
// @Description Get the contents of a single file from the main version of a project post
// @Tags 		project-posts
// @Param		projectPostID	path		string				true	"Project Post ID"
// @Param		filepath	path		string				true	"Filepath"
// @Produce		application/octet-stream
// @Success 	200		{object}	[]byte
// @Failure		404
// @Failure		500
// @Router 		/project-posts/{projectPostID}/file/{filepath}	[get]
func (projectPostController *ProjectPostController) GetMainFileFromProject(c *gin.Context) {
	// extract projectPostID
	projectPostIDStr := c.Param("projectPostID")
	projectPostID, err := strconv.ParseUint(projectPostIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid project post ID, cannot interpret as integer, id=%v ", projectPostIDStr)})

		return
	}

	// get project post for the post id
	projectPost, err := projectPostController.ProjectPostService.GetProjectPost(uint(projectPostID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no such project post found with id %v", projectPostID)})

		return
	}

	relFilepath := c.Param("filepath")
	absFilepath, err := projectPostController.PostService.GetMainFileFromProject(projectPost.PostID, relFilepath)

	// if files doesnt exist return 404 not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// get the file info
	fileContentType, err1 := mimetype.DetectFile(absFilepath)
	fileData, err2 := os.Open(absFilepath)
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
