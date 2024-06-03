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
	PostService interfaces.PostService
}

// GetPost godoc
// @Summary 	Get post
// @Description Get a post by post ID
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.PostDTO
// @Failure		400 		{object}
// @Failure		404 		{object}
// @Failure		500 		{object}
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
	post, err := postController.PostService.GetPost(postID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot get post because no post with this ID exists"})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, post)
}

// CreatePost godoc
// @Summary 	Create new post
// @Description Create a new question or discussion post
// @Accept  	json
// @Param		form	body	forms.PostCreationForm	true	"Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.PostDTO
// @Failure		400 	{object}
// @Failure		500 	{object}
// @Router 		/posts 		[post]
func (postController *PostController) CreatePost(c *gin.Context) {
	// extract post
	form := forms.PostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind PostCreationForm from request body"})

		return
	}

	// Create and add post to database here. For now just do this to test.
	post := postController.PostService.CreatePost(&form)

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &post)
}

// UpdatePost godoc
// @Summary 	Update post
// @Description Update any number of the aspects of a question or discussion post
// @Accept  	json
// @Param		post	body		models.PostDTO		true	"Updated Post"
// @Produce		json
// @Success 	200
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500 		{object}
// @Router 		/posts 		[put]
func (postController *PostController) UpdatePost(c *gin.Context) {
	// extract post
	updatedPost := models.Post{}
	err := c.BindJSON(&updatedPost)

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
// @Accept  	json
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router 		/posts/{postID} 		[delete]
func (postController *PostController) DeletePost(_ *gin.Context) {
	// delete method goes here
}

// CreatePostFromGithub godoc
// @Summary 	Create new post with the version imported from github
// @Description Create a new question or discussion post
// @Description Creates a post in the same way as CreatePost
// @Description However, the post files are imported from the given Github repository
// @Accept  	json
// @Param		form	body	forms.PostCreationForm	true	"Post Creation Form"
// @Param		url		query	string					true	"Github repository url"
// @Produce		json
// @Success 	200 	{object} 	models.PostDTO
// @Failure		400 	{object}
// @Failure		500 	{object}
// @Failure 	502 	{object}
// @Router 		/posts/from-github 		[post]
func (postController *PostController) CreatePostFromGithub(_ *gin.Context) {

}

// AddPostReport godoc
// @Summary 	Add a new report to a post
// @Description Create a new report for a post
// @Accept  	json
// @Param		form	body	forms.ReportCreationForm	true	"Report Creation Form"
// @Param		postID	path	string						true	"Post ID"
// @Produce		json
// @Success 	200 	{object} 	models.ReportDTO
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500 	{object}
// @Router 		/posts/{postID}/reports 		[post]
func (postController *PostController) AddPostReport(_ *gin.Context) {

}

// GetPostReports godoc
// @Summary		Get all reports of this post
// @Description	Get all reports that have been added to this post
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		postID		path		string			true	"Post ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.ReportDTO
// @Failure		400 	{object}
// @Failure		404 	{object}
// @Failure		500		{object}
// @Router 		/posts/{postID}/reports 		[get]
func (postController *PostController) GetPostReports(_ *gin.Context) {
	// TODO: make paginated
}
