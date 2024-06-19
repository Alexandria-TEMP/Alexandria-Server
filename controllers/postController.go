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
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Failure		500			{object} 	utils.HTTPError
// @Router 		/posts/{postID}	[get]
func (postController *PostController) GetPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret '%s' as integer: %s", postIDStr, err)})

		return
	}

	// retrieve post
	post, err := postController.PostService.GetPost(uint(postID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get post: %s", err)})

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
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.PostCreationForm	true	"Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.PostDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts 		[post]
func (postController *PostController) CreatePost(c *gin.Context) {
	form := forms.PostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind PostCreationForm from request body: %s", err)})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	// get member
	member, exists := c.Get("currentMember")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get logged in user"})

		return
	}

	post, err := postController.PostService.CreatePost(&form, member.(*models.Member))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create post: %s", err)})

		return
	}

	// response
	c.JSON(http.StatusOK, &post)
}

// DeletePost godoc
// @Summary 	Delete a post
// @Description Delete a post with given ID from database
// @Tags 		posts
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts/{postID} 		[delete]
func (postController *PostController) DeletePost(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// CreatePostFromGithub godoc
// @Summary 	Create new post with the version imported from github
// @Description Create a new question or discussion post
// @Description Creates a post in the same way as CreatePost
// @Description However, the post files are imported from the given Github repository
// @Tags 		posts
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.PostCreationForm	true	"Post Creation Form"
// @Param		url		query	string					true	"Github repository url"
// @Produce		json
// @Success 	200 	{object} 	models.PostDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Failure 	502		{object} 	utils.HTTPError
// @Router 		/posts/from-github 		[post]
func (postController *PostController) CreatePostFromGithub(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// AddPostReport godoc
// @Summary 	Add a new report to a post
// @Description Create a new report for a post
// @Tags 		posts
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.ReportCreationForm	true	"Report Creation Form"
// @Param		postID	path	string						true	"Post ID"
// @Produce		json
// @Success 	200 	{object} 	models.ReportDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/reports 		[post]
func (postController *PostController) AddPostReport(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetPostReports godoc
// @Summary		Get all reports of this post
// @Description	Get all reports that have been added to this post
// @Tags 		posts
// @Accept 		json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/reports 		[get]
func (postController *PostController) GetPostReports(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetCollaborator godoc
// @Summary 	Get a post collaborator by ID
// @Description	Get a post collaborator by ID, a member who has collaborated on a post
// @Tags		posts
// @Accept  	json
// @Param		collaboratorID	path	string	true	"Collaborator ID"
// @Produce		json
// @Success 	200 		{object}	models.PostCollaboratorDTO
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Failure		500			{object} 	utils.HTTPError
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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/posts/reports/{reportID}				[get]
func (postController *PostController) GetPostReport(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// UploadPost
// @Summary 	Upload a new project version to a branch
// @Description Upload a zipped quarto project to a post. This is the main version of the post, as there are no other versions.
// @Description Specifically, this zip should contain all of the contents of the project at its root, not in a subdirectory.
// @Tags 		posts
// @Accept  	multipart/form-data
// @Param 		Authorization header string true "Access Token"
// @Param		postID			path		string			true	"Post ID"
// @Param		file			formData	file			true	"Repository to create"
// @Produce		application/json
// @Success 	200
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/upload		[post]
func (postController *PostController) UploadPost(c *gin.Context) {
	// extract file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("no file found: %s", err)})

		return
	}

	// extract post id
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID '%s', cannot interpret as integer: %s", postIDStr, err)})

		return
	}

	// Create commit on branch with new files
	err = postController.PostService.UploadPost(c, file, uint(postID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/render	[get]
func (postController *PostController) GetMainRender(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID '%s', cannot interpret as integer: %s", postIDStr, err)})

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
	c.Status(http.StatusOK)
}

// GetMainProject godoc specs are subject to change
// @Summary 	Get the main repository of a post
// @Description Get the entire zipped main repository of a post
// @Tags 		posts
// @Param		postID	path		string				true	"Post ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/repository	[get]
func (postController *PostController) GetMainProject(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID '%s', cannot interpret as integer: %s", postIDStr, err)})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/tree		[get]
func (postController *PostController) GetMainFiletree(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID '%s', cannot interpret as integer: %s", postIDStr, err)})

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
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
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
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/posts/{postID}/project-post	[get]
func (postController *PostController) GetProjectPostIfExists(c *gin.Context) {
	// Note: this endpoint kind of goes against the data model's design philosophy, and is quite a hacky fix.
	// TODO reconsider how composition of posts and project posts is implemented & integrated.
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
	c.JSON(http.StatusOK, projectPost.ID)
}

// GetAllPostCollaborators godoc
// @Summary 	Get all post collaborators of a post
// @Description Returns all post collaborators of the post with the given ID
// @Tags 		posts
// @Param		postID	path		string		true	"Post ID"
// @Produce		application/json
// @Success 	200		{array}		models.PostCollaboratorDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router		/posts/collaborators/all/{postID}		[get]
func (postController *PostController) GetAllPostCollaborators(c *gin.Context) {
	// Get post ID from path param
	postIDString := c.Param("postID")

	postID, err := strconv.ParseUint(postIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse post ID '%s' as unsigned integer: %s", postIDString, err)})

		return
	}

	// Get the post itself
	post, err := postController.PostService.GetPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get post with ID %d: %s", postID, err)})

		return
	}

	postCollaborators := post.Collaborators

	// Turn each post collaborator into a DTO
	postCollaboratorDTOs := make([]*models.PostCollaboratorDTO, len(postCollaborators))

	for i, postCollaborator := range postCollaborators {
		postCollaboratorDTO := postCollaborator.IntoDTO()
		postCollaboratorDTOs[i] = &postCollaboratorDTO
	}

	c.JSON(http.StatusOK, postCollaboratorDTOs)
}
