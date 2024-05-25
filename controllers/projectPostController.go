package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v2

type ProjectPostController struct {
}

// GetProjectPost godoc
// @Summary 	Get project post
// @Description Get a project post by post ID
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.ProjectPostDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		410 		{object} 	utils.HTTPError
// @Router 		/project-post/{postID}	[get]
func (postController *PostController) GetProjectPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid post ID, cannot interpret as integer, id=%s ", postIDStr))

		return
	}

	post, err := postController.PostService.GetProjectPost(uint64(postID))

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
// @Param		form	body		forms.ProjectPostCreationForm	true	"Project Post Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.ProjectPostDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Router 		/projectPost 		[post]
func (postController *PostController) CreateProjectPost(c *gin.Context) {
	// extract post
	form := forms.ProjectPostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind ProjectPostCreationForm from request body"))

		return
	}

	// Create and add post to database here. For now just do this to test.
	post := postController.PostService.CreateProjectPost(&form)

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &post)
}

// UpdateProjectPost godoc
// @Summary 	Update project post
// @Description Update any number of the aspects of project post
// @Accept  	json
// @Param		post	body		models.ProjectPostDTO		true	"Updated Project Post"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		410 	{object} 	utils.HTTPError
// @Router 		/ 		[put]
func (postController *PostController) UpdateProjectPost(c *gin.Context) {
	// extract post
	updatedProjectPost := models.ProjectPost{}
	err := c.BindJSON(&updatedProjectPost)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind updated ProjectPost from request body"))

		return
	}

	// Update and add post to database here. For now just do this to test.
	err = postController.PostService.UpdateProjectPost(&updatedProjectPost)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusGone, errors.New("cannot update post because no ProjectPost with this ID exists"))

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}