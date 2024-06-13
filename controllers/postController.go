package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type PostController struct {
	PostService             interfaces.PostService
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
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
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
	c.Header("Content-Type", "application/json")
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
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

	post, err := postController.PostService.CreatePost(&form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create post: %s", err)})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/posts 		[put]
func (postController *PostController) UpdatePost(c *gin.Context) {
	// extract post
	updatedPost := models.Post{}
	err := c.BindJSON(&updatedPost)

	// TODO convert from Post DTO to updated Post data

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind updated Post from request body: %s", err)})

		return
	}

	// Update and add post to database here. For now just do this to test.
	err = postController.PostService.UpdatePost(&updatedPost)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot update post because no post with this ID exists: %s", err)})

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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Failure 	502 	{object}	utils.HTTPError
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
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
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
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
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get post collaborator: %s", err)})

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
// @Failure		400		{object}	utils.HTTPError
// @Failure		404		{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/posts/reports/{reportID}				[get]
func (postController *PostController) GetPostReport(_ *gin.Context) {
	// TODO implement
}
