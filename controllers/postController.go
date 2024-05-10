package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v1

type PostController struct {
	PostService services.PostService
}

// GetPost godoc
// @Summary Get post
// @Description Get post by post ID
// @Accept  	json
// @Param		postID	path	string	true	"postID"
// @Produce		json
// @Success 	200 {object} models.Post
// @Failure		400 {object} utils.HTTPError
// @Router 		/post/{postID} [get]
func (postController *PostController) GetPost(c *gin.Context) {
	// extract postID
	postIDStr := c.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid post ID, cannot interpret as integer, id=%s ", postIDStr))

		return
	}

	// extract versionID
	versionIDStr := c.Param("versionID")
	versionID, err := strconv.ParseInt(versionIDStr, 10, 64)

	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid version ID, cannot interpret as integer, id=%s ", versionIDStr))

		return
	}

	fmt.Printf("GET /post/%v/%v\n", versionID, postID)

	// get post from database here. For now just do this to test.
	post := new(models.Post)

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, post)
}

// // CreatePost godoc
// func (postController *PostController) CreatePost(c *gin.Context) {

// }
