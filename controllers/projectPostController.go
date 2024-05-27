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
	//TODO: change to project post service
	ProjectPostService interfaces.PostService
}

// GetProjectPost godoc
// @Summary 	Get project post
// @Description Get a project post by ID
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
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid post ID, cannot interpret as integer, id=%s ", postIDStr))

		return
	}

	post, err := projectPostController.ProjectPostService.GetProjectPost(uint64(postID))

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusGone, errors.New("cannot get project post because no post with this ID exists"))

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, post)
}

// CreateProjectPost godoc
// @Summary 	Create new project post
// @Description Create a new project post
// @Accept  	json
// @Param		form			body		forms.ProjectPostCreationForm	true	"Project Post Creation Form"
// @Param 		parentPostID	query		string							false	"Parent post ID"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/project-posts		[post]
func (projectPostController *ProjectPostController) CreateProjectPost(c *gin.Context) {
	// extract post
	form := forms.ProjectPostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind ProjectPostCreationForm from request body"))

		return
	}

	// Create and add post to database here. For now just do this to test.
	post := projectPostController.ProjectPostService.CreateProjectPost(&form)

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &post)
}

// UpdateProjectPost godoc
// @Summary 	Update project post
// @Description Update any number of the aspects of a project post
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
// @Summary Returns all discussions associated with the project post
// @Description Returns all discussions on this project post and all of it's merge requests
// @Description Endpoint is offset-paginated
// @Accept  	json
// @Param		postID		path		string			true	"post ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/project-posts/{postID}/all-discussions 	[get]
func (projectPostController *ProjectPostController) GetProjectPostDiscussions(_ *gin.Context) {

}

// GetProjectPostOpenMergeRequests godoc
// @Summary		Get all open merge requests of a project post
// @Description	Get all open merge requests associated with the given project post
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		postID		path		string			true	"post ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.MergeRequestDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/project-posts/{postID}/open-merge-requests 		[get]
func (projectPostController *ProjectPostController) GetProjectPostOpenMergeRequests(_ *gin.Context) {
	// return all the merge requests associated with this project post that are open
	// TODO: make endpoint paginated
}

// GetProjectPostClosedMergeRequests godoc
// @Summary		Get all closed merge requests of a project post
// @Description	Get all closed merge requests associated with the given project post
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		postID		path		string			true	"post ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.MergeRequestDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/project-posts/{postID}/closed-merge-requests 		[get]
func (projectPostController *ProjectPostController) GetProjectPostClosedMergeRequests(_ *gin.Context) {
	// return all the merge requests associated with this project post that are closed
	// TODO: make endpoint paginated
}
