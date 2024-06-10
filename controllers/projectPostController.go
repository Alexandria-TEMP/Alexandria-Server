package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v2

type ProjectPostController struct {
	ProjectPostService interfaces.ProjectPostService
}

// GetProjectPost godoc
// @Summary 	Get project post
// @Description Get a project post by ID
// @Tags 		project-posts
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.ProjectPostDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
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
// @Description Create a new project post
// @Tags 		project-posts
// @Accept  	json
// @Param		form	body		forms.ProjectPostCreationForm	true	"Project Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
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

	projectPost, err := projectPostController.ProjectPostService.CreateProjectPost(&form)
	if err != nil {
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
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/project-posts 		[put]
func (projectPostController *ProjectPostController) UpdateProjectPost(c *gin.Context) {
	// extract post
	updatedProjectPost := models.ProjectPost{}
	err := c.BindJSON(&updatedProjectPost)

	// TODO convert from project post DTO to updated project post

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind updated ProjectPost from request body"))

		return
	}

	// Update and add post to database here. For now just do this to test.
	err = projectPostController.ProjectPostService.UpdateProjectPost(&updatedProjectPost)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusGone, errors.New("cannot update post because no ProjectPost with this ID exists"))

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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Failure 	502 	{object}	utils.HTTPError
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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/project-posts/{postID}/all-discussions 	[get]
func (projectPostController *ProjectPostController) GetProjectPostDiscussions(_ *gin.Context) {
	// TODO implement
}

// GetProjectPostMRsByStatus godoc
// @Summary 	Returns MR IDs grouped by each MR status
// @Description Returns all MR IDs of this project post, grouped by each MR's review status
// @Tags		project-posts
// @Accept		json
// @Param		postID	path	string	true	"post ID"
// @Produce		json
// @Success		200		{object}	forms.GroupedBranchForm
// @Failure		400		{object}	utils.HTTPError
// @Failure		404		{object}	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/project-posts/{postID}/branches-by-status	[get]
func (projectPostController *ProjectPostController) GetProjectPostMRsByStatus(_ *gin.Context) {
	// TODO implement
}
