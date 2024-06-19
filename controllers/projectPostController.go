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

type ProjectPostController struct {
	ProjectPostService         interfaces.ProjectPostService
	DiscussionContainerService interfaces.DiscussionContainerService
	PostService                interfaces.PostService
	RenderService              interfaces.RenderService
}

// GetProjectPost godoc
// @Summary 	Get project post
// @Description Get a project post by ID
// @Tags 		project-posts
// @Accept  	json
// @Param		projectPostID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 	{object}	models.ProjectPostDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/project-posts/{projectPostID}	[get]
func (projectPostController *ProjectPostController) GetProjectPost(c *gin.Context) {
	// extract projectPostID
	projectPostIDStr := c.Param("projectPostID")
	projectPostID, err := strconv.ParseUint(projectPostIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not interpret ID %s as unsigned integer: %s", projectPostIDStr, err)})

		return
	}

	projectPost, err := projectPostController.ProjectPostService.GetProjectPost(uint(projectPostID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not get project post: %v", err.Error())})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, projectPost.IntoDTO())
}

// CreateProjectPost godoc
// @Summary 	Create new project post
// @Description Create a new project post with a single open branch. Upload to this branch in order to have your post reviewed.
// @Tags 		project-posts
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body		forms.ProjectPostCreationForm	true	"Project Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/project-posts		[post]
func (projectPostController *ProjectPostController) CreateProjectPost(c *gin.Context) {
	form := forms.ProjectPostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid project post creation form: %v", err.Error())})

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

	projectPost, err404, err500 := projectPostController.ProjectPostService.CreateProjectPost(&form, member.(*models.Member))

	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("not found: %v", err404.Error())})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("internal server error: %v", err500.Error())})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, projectPost.IntoDTO())
}

// DeleteProjectPost godoc
// @Summary 	Delete a project post
// @Description Delete a project post with given ID from database
// @Tags 		project-posts
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		projectPostID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/project-posts/{projectPostID} 		[delete]
func (projectPostController *ProjectPostController) DeleteProjectPost(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// CreateProjectPostFromGithub godoc
// @Summary 	Create new project post with the version imported from github
// @Description Create a new project post
// @Description Creates a project post in the same way as CreateProjectPost
// @Description However, the post files are imported from the given Github repository
// @Tags 		project-posts
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.ProjectPostCreationForm	true	"Post Creation Form"
// @Param		url		query	string							true	"Github repository url"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Failure 	502		{object} 	utils.HTTPError
// @Router 		/project-posts/from-github 		[post]
func (projectPostController *ProjectPostController) CreateProjectPostFromGithub(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetProjectPostDiscussionContainers godoc
// @Summary Returns all discussion container IDs associated with the project post
// @Description Returns all discussion container IDs on this project post over all its previous merged versions, instead of only the current version
// @Tags 		project-posts
// @Accept  	json
// @Param		projectPostID		path		string			true	"project post ID"
// @Produce		json
// @Success 	200		{object}	models.DiscussionContainerProjectHistoryDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/project-posts/{projectPostID}/all-discussion-containers 	[get]
func (projectPostController *ProjectPostController) GetProjectPostDiscussionContainers(c *gin.Context) {
	// Get project post ID from path
	projectPostIDString := c.Param("projectPostID")

	projectPostID, err := strconv.ParseUint(projectPostIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse project post ID '%s' as unsigned integer: %s", projectPostIDString, err)})

		return
	}

	// Get the discussion container history!
	discussionContainerHistory, err := projectPostController.ProjectPostService.GetDiscussionContainersFromMergeHistory(uint(projectPostID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get discussion containers of project post: %s", err)})

		return
	}

	c.JSON(http.StatusOK, discussionContainerHistory)
}

// GetProjectPostBranchesByStatus godoc
// @Summary 	Get branch IDs by review status
// @Description Returns all branch IDs of this project post, grouped by each branch's review status
// @Tags		project-posts
// @Accept		json
// @Param		projectPostID	path	string	true	"project post ID"
// @Produce		json
// @Success		200		{object}	models.BranchesGroupedByReviewStatusDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/project-posts/{projectPostID}/branches-by-status	[get]
func (projectPostController *ProjectPostController) GetProjectPostBranchesByStatus(c *gin.Context) {
	// Get project post ID from path
	projectPostIDString := c.Param("projectPostID")

	projectPostID, err := strconv.ParseUint(projectPostIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse project post ID '%s' as unsigned integer: %s", projectPostIDString, err)})

		return
	}

	branchesGroupedByStatus, err := projectPostController.ProjectPostService.GetBranchesGroupedByReviewStatus(uint(projectPostID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get branches by status: %s", err)})

		return
	}

	c.JSON(http.StatusOK, branchesGroupedByStatus)
}
