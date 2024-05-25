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
// @Success 	200 		{object}	models.PostDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		410 		{object} 	utils.HTTPError
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
	post, err := postController.PostService.GetPost(postID)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusGone, errors.New("cannot get post because no post with this ID exists"))

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

// UpdatePost godoc
// @Summary 	Update post
// @Description Update any number of the aspects of a question or discussion post
// @Accept  	json
// @Param		post	body		models.PostDTO		true	"Updated Post"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		410 	{object} 	utils.HTTPError
// @Router 		/ 		[put]
func (postController *PostController) UpdatePost(c *gin.Context) {
	// extract post
	updatedPost := models.Post{}
	err := c.BindJSON(&updatedPost)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind updated Post from request body"))

		return
	}

	// Update and add post to database here. For now just do this to test.
	err = postController.PostService.UpdatePost(&updatedPost)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusGone, errors.New("cannot update post because no post with this ID exists"))

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}

