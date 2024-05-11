package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v1

type PostController struct {
	PostService interfaces.PostService
}

// GetPost godoc
// @Summary 	Get post
// @Description Get a post by post ID
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.Post
// @Failure		400 		{object} 	utils.HTTPError
// @Router 		/post/{postID}	[get]
func (postController *PostController) GetPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid post ID, cannot interpret as integer, id=%s ", postIDStr))

		return
	}

	// retrieve post
	post := postController.PostService.GetPost(postID)

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
// @Success 	200 	{object} 	models.Post
// @Failure		400 	{object} 	utils.HTTPError
// @Router 		/post 		[post]
func (postController *PostController) CreatePost(c *gin.Context) {
	// extract post
	form := forms.PostCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind PostCreationForm from request body"))

		return
	}

	// Create and add post to database here. For now just do this to test.
	post := postController.PostService.CreatePost(&form)

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &post)
}

// GetProjectPost godoc
// @Summary 	Get project post
// @Description Get a project post by post ID
// @Accept  	json
// @Param		postID		path		string			true	"Post ID"
// @Produce		json
// @Success 	200 		{object}	models.ProjectPost
// @Failure		400 		{object} 	utils.HTTPError
// @Router 		/projectPost/{postID}	[get]
func (postController *PostController) GetProjectPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid post ID, cannot interpret as integer, id=%s ", postIDStr))

		return
	}

	post := postController.PostService.GetProjectPost(uint64(postID))

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
// @Success 	200 	{object} 	models.ProjectPost
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
