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

type PostController struct {
	PostService             interfaces.PostService
	RenderService           interfaces.RenderService
	PostCollaboratorService interfaces.PostCollaboratorService
}

// GetPost godoc
// @Summary 	Get post by ID
// @Description Get a post by ID
// @Tags 		posts
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.PostDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID}	[get]
func (postController *PostController) GetPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret as integer, id=%s ", postIDStr)})

		return
	}

	// retrieve post
	post, err := postController.PostService.GetPost(uint(postID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot get post because no post with this ID exists"})

		return
	}

	// response
	c.JSON(http.StatusOK, post)
}

// CreatePost godoc
// @Summary 	Create new post
// @Description Create a new question or discussion post. Cannot be a project post.
// @Tags 		posts
// @Accept  	json
// @Param		form	body	forms.PostCreationForm	true	"Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.PostDTO
// @Failure		400
// @Failure		500
// @Router 		/posts 		[post]
func (postController *PostController) CreatePost(c *gin.Context) {
	form := forms.PostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind PostCreationForm from request body"})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	post, err := postController.PostService.CreatePost(&form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create post, reason: %v", err.Error())})

		return
	}

	// response
	c.JSON(http.StatusOK, &post)
}

// UpdatePost godoc
// @Summary 	Update post
// @Description Update any number of aspects of a question or discussion post
// @Tags 		posts
// @Accept  	json
// @Param		post	body		models.PostDTO		true	"Updated Post"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts 		[put]
func (postController *PostController) UpdatePost(c *gin.Context) {
	// extract post
	updatedPost := models.Post{}
	err := c.BindJSON(&updatedPost)

	// TODO convert from Post DTO to updated Post data

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind updated Post from request body"})

		return
	}

	// Update and add post to database here. For now just do this to test.
	err = postController.PostService.UpdatePost(&updatedPost)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot update post because no post with this ID exists"})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}

// DeletePost godoc
// @Summary 	Delete a post
// @Description Delete a post with given ID from database
// @Tags 		posts
// @Accept  	json
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID} 		[delete]
func (postController *PostController) DeletePost(_ *gin.Context) {
	// delete method goes here
}

// CreatePostFromGithub godoc
// @Summary 	Create new post with the version imported from github
// @Description Create a new question or discussion post
// @Description Creates a post in the same way as CreatePost
// @Description However, the post files are imported from the given Github repository
// @Tags 		posts
// @Accept  	json
// @Param		form	body	forms.PostCreationForm	true	"Post Creation Form"
// @Param		url		query	string					true	"Github repository url"
// @Produce		json
// @Success 	200 	{object} 	models.PostDTO
// @Failure		400
// @Failure		500
// @Failure 	502
// @Router 		/posts/from-github 		[post]
func (postController *PostController) CreatePostFromGithub(_ *gin.Context) {

}

// AddPostReport godoc
// @Summary 	Add a new report to a post
// @Description Create a new report for a post
// @Tags 		posts
// @Accept  	json
// @Param		form	body	forms.ReportCreationForm	true	"Report Creation Form"
// @Param		postID	path	string						true	"Post ID"
// @Produce		json
// @Success 	200 	{object} 	models.ReportDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID}/reports 		[post]
func (postController *PostController) AddPostReport(_ *gin.Context) {

}

// GetPostReports godoc
// @Summary		Get all reports of this post
// @Description	Get all reports that have been added to this post
// @Tags 		posts
// @Accept 		json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID}/reports 		[get]
func (postController *PostController) GetPostReports(_ *gin.Context) {
	// TODO implement
}

// GetCollaborator godoc
// @Summary 	Get a post collaborator by ID
// @Description	Get a post collaborator by ID, a member who has collaborated on a post
// @Tags		posts
// @Accept  	json
// @Param		collaboratorID	path	string	true	"Collaborator ID"
// @Produce		json
// @Success 	200 		{object}	models.PostCollaboratorDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts/collaborators/{collaboratorID}	[get]
func (postController *PostController) GetPostCollaborator(c *gin.Context) {
	idString := c.Param("collaboratorID")

	// Parse path parameter into an integer
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse ID '%s' as unsigned integer: %s", idString, err)})

		return
	}

	// Fetch the post collaborator by ID
	postCollaborator, err := postController.PostCollaboratorService.GetPostCollaborator(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get post collaborator: %v", err.Error())})

		return
	}

	c.JSON(http.StatusOK, postCollaborator)
}

// GetPostReport godoc
// @Summary		Gets a post report by ID
// @Description	Gets a post report by its ID
// @Tags		posts
// @Param		reportID	path	string	true	"Report ID"
// @Produce		json
// @Success		200		{object}	reports.PostReportDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/posts/reports/{reportID}				[get]
func (postController *PostController) GetPostReport(_ *gin.Context) {
	// TODO implement
}

// UploadPost
// @Summary 	Upload a new project version to a branch
// @Description Upload a zipped quarto project to a post. This is the main version of the post, as there are no other versions.
// @Description Specifically, this zip should contain all of the contents of the project at its root, not in a subdirectory.
// @Tags 		posts
// @Accept  	multipart/form-data
// @Param		postID			path		string			true	"Post ID"
// @Param		file			formData	file			true	"Repository to create"
// @Produce		application/json
// @Success 	200
// @Failure		400
// @Failure		500
// @Router 		/posts/{postID}/upload		[post]
func (postController *PostController) UploadPost(c *gin.Context) {
	// extract file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file found"})

		return
	}

	// extract post id
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%v ", postIDStr)})

		return
	}

	// Create commit on branch with new files
	err = postController.PostService.UploadPost(c, file, uint(postID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Status(http.StatusOK)
}

// GetMainRender
// @Summary 	Get the main render of a post
// @Description Get the main render of the repository underlying a post if it exists and has been rendered successfully
// @Tags 		posts
// @Param		postID		path		string				true	"Post ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Success		202		{object}	[]byte
// @Failure		400
// @Failure		404
// @Router 		/posts/{postID}/render	[get]
func (postController *PostController) GetMainRender(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret as integer, id=%v ", postIDStr)})

		return
	}

	// get render filepath
	filePath, err202, err404 := postController.RenderService.GetMainRenderFile(uint(postID))

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
// @Summary 	Get the main repository of a post
// @Description Get the entire zipped main repository of a post
// @Tags 		posts
// @Param		postID	path		string				true	"Post ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		400
// @Failure		404
// @Router 		/posts/{postID}/repository	[get]
func (postController *PostController) GetMainProject(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret as integer, id=%v ", postIDStr)})

		return
	}

	// get repository filepath
	filePath, err := postController.PostService.GetMainProject(uint(postID))

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
// @Summary 	Get the filetree of a post
// @Description Get the filetree of a the main version of a post, together with the size of the file in bytes.
// @Description Directories have a size of -1.
// @Tags 		posts
// @Param		postID	path		string				true	"Post ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID}/tree		[get]
func (postController *PostController) GetMainFiletree(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret as integer, id=%v ", postIDStr)})

		return
	}

	fileTree, err404, err500 := postController.PostService.GetMainFiletree(uint(postID))

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
// @Summary 	Get a file from a post
// @Description Get the contents of a single file from the main version of a post
// @Tags 		posts
// @Param		postID	path		string				true	"Post ID"
// @Param		filepath	path		string				true	"Filepath"
// @Produce		application/octet-stream
// @Success 	200		{object}	[]byte
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID}/file/{filepath}	[get]
func (postController *PostController) GetMainFileFromProject(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret as integer, id=%v ", postIDStr)})

		return
	}

	relFilepath := c.Param("filepath")
	absFilepath, err := postController.PostService.GetMainFileFromProject(uint(postID), relFilepath)

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

// GetProjectPostIfExists godoc
// @Summary 	Get Project Post of Post
// @Description Get the Project Post ID that encapsulates a Post, if this Project Post exists
// @Tags 		posts
// @Param		postID		path		string				true	"Post ID"
// @Produce		application/json
// @Success 	200		{object}	uint
// @Failure		404
// @Failure		500
// @Router 		/posts/{postID}/project-post	[get]
func (postController *PostController) GetProjectPostIfExists(c *gin.Context) {
	// Get post ID from path
	postIDString := c.Param("postID")

	postID, err := strconv.ParseUint(postIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse post ID '%s' as unsigned integer: %s", postIDString, err)})

		return
	}

	// Get the post's project post
	projectPost, err := postController.PostService.GetProjectPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get project post of post with ID %d: %s", postID, err)})

		return
	}

	// Return the project post's ID
	c.JSON(http.StatusOK, gin.H{"projectPostID": projectPost.ID})

	// Note: this endpoint kind of goes against the data model's design philosophy, and is quite a hacky fix.
	// TODO reconsider how composition of posts and project posts is implemented & integrated.
}
