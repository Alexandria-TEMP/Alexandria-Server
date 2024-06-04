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

type MemberController struct {
	MemberService interfaces.MemberService
}

// GetMember godoc
// @Summary 	Get member from database
// @Description Get a member by user ID
// @Accept  	json
// @Param		userID		path		string			true	"user ID"
// @Produce		json
// @Success 	200 		{object}	models.MemberDTO
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}	[get]
func (memberController *MemberController) GetMember(c *gin.Context) {
	// extract the id of the member
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	// if this caused an error, print it
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid user ID, cannot interpret as integer, id=%s ", userIDStr)})

		return
	}
	// get the user through the service
	member, err := memberController.MemberService.GetMember(uint(userID))

	// if there was an error, print it and return
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot get member because no user with this ID exists"})

		return
	}

	// if correct response send the member back
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, member)
}

// CreateMember godoc
// @Summary 	Create a new member
// @Description Create a new member from the given fields
// @Accept  	json
// @Param		form	body	forms.MemberCreationForm	true	"Member Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.MemberDTO
// @Failure		400
// @Failure		500
// @Router 		/members 		[post]
func (memberController *MemberController) CreateMember(c *gin.Context) {
	// get the member
	form := forms.MemberCreationForm{}
	// bind the fields of the param to the JSON of the model
	err := c.BindJSON(&form)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind userCreationForm from request body"})

		return
	}

	// create and add to database(not done yet) through the memberService
	member := memberController.MemberService.CreateMember(&form)

	// send back a positive response with the created member
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &member)
}

// UpdateMember godoc
// @Summary 	Update a member
// @Description Update the fields of a member
// @Accept  	json
// @Param		member	body		models.Member		true	"Updated member"
// @Produce		json
// @Success 	200
// @Failure		404
// @Failure		500
// @Router 		/members 		[put]
func (memberController *MemberController) UpdateMember(c *gin.Context) {
	// get the new member object
	updatedMember := models.Member{}
	err := c.BindJSON(&updatedMember)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind updated user from request body"})

		return
	}

	// update and add the member to the database
	err = memberController.MemberService.UpdateMember(&updatedMember)

	// check for errors again
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cannot update user because no user with this ID exists"})

		return
	}

	// send back a positive response if member updated successfully
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}

// DeleteMember godoc
// @Summary 	Delete a member
// @Description Delete a member with given ID from database
// @Accept  	json
// @Param		userID		path		string			true	"user ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID} 		[delete]
func (memberController *MemberController) DeleteMember(_ *gin.Context) {
	// delete method goes here
}

// GetMemberPosts godoc
// @Summary		Get all posts of this member
// @Description	Get all posts that this member is a collaborator of
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		userID		path		string			true	"user ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.PostDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}/posts 		[get]
func (memberController *MemberController) GetMemberPosts(_ *gin.Context) {
	// return all the posts
	// that this member is a collaborator/author of
	// TODO: make endpoint paginated
}

// GetMemberProjectPosts godoc
// @Summary		Get all project posts of this member
// @Description	Get all project posts that this member is a collaborator of
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		userID		path		string			true	"user ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.ProjectPostDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}/project-posts 		[get]
func (memberController *MemberController) GetMemberProjectPosts(_ *gin.Context) {
	// return all the project posts
	// that this member is a collaborator/author of
	// TODO: make endpoint paginated
}

// GetMemberBranches godoc
// @Summary		Get all branches of this member
// @Description	Get all branches that this member is a collaborator of
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		userID		path		string			true	"user ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.BranchDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}/branches 		[get]
func (memberController *MemberController) GetMemberBranches(_ *gin.Context) {
	// return all the branches
	// that this member is a collaborator/author of
	// TODO: make endpoint paginated
}

// GetMemberDiscussions godoc
// @Summary		Get all branches of this member
// @Description	Get all branches that this member is a collaborator of
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		userID		path		string			true	"user ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}/discussions		[get]
func (memberController *MemberController) GetMemberDiscussions(_ *gin.Context) {
	// returns all the discussions this member is a part of
	// TODO: make paginated
}

// AddMemberSavedPost godoc
// @Summary 	Adds new saved post
// @Description Adds a post to the saved posts of a member
// @Accept  	json
// @Param		userID		path		string			true	"user ID"
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		500
// @Router 		/members/{userID}/saved-posts/{postID} 		[post]
func (memberController *MemberController) AddMemberSavedPost(_ *gin.Context) {

}

// AddMemberSavedProjectPost godoc
// @Summary 	Adds new saved project post
// @Description Adds a project post to the saved project posts of a member
// @Accept  	json
// @Param		userID		path		string			true	"user ID"
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		500
// @Router 		/members/{userID}/saved-project-posts/{postID} 		[post]
func (memberController *MemberController) AddMemberSavedProjectPost(_ *gin.Context) {

}

// GetMemberSavedPosts godoc
// @Summary		Get all saved posts of this member
// @Description	Get all posts that this member has saved
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		userID		path		string			true	"user ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.PostDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}/saved-posts 		[get]
func (memberController *MemberController) GetMemberSavedPosts(_ *gin.Context) {
	// return all saved posts of this member
	// TODO: make endpoint paginated
}

// GetMemberProjectPosts godoc
// @Summary		Get all saved project posts of this member
// @Description	Get all project posts that this member has saved
// @Description Endpoint is offset-paginated
// @Accept 		json
// @Param		userID		path		string			true	"user ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		json
// @Success 	200		{array}		models.ProjectPostDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/{userID}/saved-project-posts 		[get]
func (memberController *MemberController) GetMemberSavedProjectPosts(_ *gin.Context) {
	// return all the project posts that this member has saved
	// TODO: make endpoint paginated
}
