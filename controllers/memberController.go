package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

// @BasePath /api/v2

type MemberController struct {
	MemberService interfaces.MemberService
	TagService    interfaces.TagService
}

// GetMember godoc
// @Summary 	Get member from database
// @Description Get a member by member ID
// @Tags 		members
// @Accept  	json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200 		{object}	models.MemberDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500			{object}	utils.HTTPError
// @Router 		/members/{memberID}	[get]
func (memberController *MemberController) GetMember(c *gin.Context) {
	// extract the id of the member
	memberIDStr := c.Param("memberID")
	initmemberID, err := strconv.ParseUint(memberIDStr, 10, 64)
	// if this caused an error, print it and return status 400: bad input
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid member ID, cannot interpret as integer, id=%s ", memberIDStr))

		return
	}

	// cast member ID as uint instead of uint64, because database only accepts those
	memberID := uint(initmemberID)

	// get the member through the service
	member, err := memberController.MemberService.GetMember(memberID)

	// move the member into a DTO format
	memberDTO := member.IntoDTO()

	// if there was an error, print it and return status 404: not found
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusNotFound, fmt.Errorf("cannot get member because no member with this ID exists, id=%d", memberID))

		return
	}

	// if correct response send the member back
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, memberDTO)
}

// CreateMember godoc
// @Summary 	Create a new member
// @Description Create a new member from the given fields
// @Tags 		members
// @Accept  	json
// @Param		form	body	forms.MemberCreationForm	true	"Member Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.MemberDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members 		[post]
func (memberController *MemberController) CreateMember(c *gin.Context) {
	// get the member
	form := forms.MemberCreationForm{}
	// bind the fields of the param to the JSON of the model
	err := c.BindJSON(&form)

	if !form.IsValid() {
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("form fails validation"))
		return
	}

	// check for errors
	// if there is an error, return a 400 bad request status
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind memberCreationForm from request body"))

		return
	}

	// get array of strings, create array of tags
	tagIDs := form.ScientificFieldTagIDs

	// call the method from the tag service
	tagArray, err := memberController.TagService.GetTagsFromUintIDs(tagIDs)

	tagContainer := tags.ScientificFieldTagContainer{
		ScientificFieldTags: tagArray,
	}

	// check for errors in the tags
	// if there is an error, return a 400 bad request status
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind tag ids from request body"))

		return
	}

	// create and add to database through the memberService
	member, err := memberController.MemberService.CreateMember(&form, &tagContainer)

	memberDTO := member.IntoDTO()

	// if the member service throws an error, return a 400 Bad request status
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("a member with this id already exists"))

		return
	}

	// send back a positive response 200 status with the created member
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, memberDTO)
}

// UpdateMember godoc
// @Summary 	Update a member
// @Description Update the fields of a member
// @Tags 		members
// @Accept  	json
// @Param		member	body		models.MemberDTO		true	"Updated member"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members 		[put]
func (memberController *MemberController) UpdateMember(c *gin.Context) {
	// get the new member object
	updatedMember := models.MemberDTO{}
	err := c.BindJSON(&updatedMember)

	// check for errors
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind updated member from request body"))

		return
	}

	// get array of strings, create array of tags
	tagIDs := updatedMember.ScientificFieldTagIDs
	// call the method from the tag service
	tagArray, err := memberController.TagService.GetTagsFromUintIDs(tagIDs)
	tagContainer := tags.ScientificFieldTagContainer{
		ScientificFieldTags: tagArray,
	}

	// check for errors in the tags
	// if there is an error, return a 400 bad request status
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusBadRequest, errors.New("cannot bind tag ids from request body"))

		return
	}

	err = memberController.MemberService.UpdateMember(&updatedMember, &tagContainer)
	// check for errors again
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusNotFound, errors.New("cannot update member because no member with this ID exists"))

		return
	}

	// send back a positive response if member updated successfully
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}

// DeleteMember godoc
// @Summary 	Delete a member
// @Description Delete a member with given ID from database
// @Tags 		members
// @Accept  	json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID} 		[delete]
func (memberController *MemberController) DeleteMember(c *gin.Context) {
	// extract the id of the member
	memberIDStr := c.Param("memberID")
	initmemberID, err1 := strconv.ParseUint(memberIDStr, 10, 64)

	// if this caused an error, print it and return status 400: bad input
	if err1 != nil {
		fmt.Println(err1)
		utils.ThrowHTTPError(c, http.StatusBadRequest, fmt.Errorf("invalid member ID, cannot interpret as integer, id=%s ", memberIDStr))

		return
	}

	// cast member ID as uint instead of uint64, because database only accepts those
	memberID := uint(initmemberID)

	// get the member through the service
	err2 := memberController.MemberService.DeleteMember(memberID)

	// if there was an error, print it and return status 404: not found
	if err2 != nil {
		fmt.Println(err2)
		utils.ThrowHTTPError(c, http.StatusNotFound, fmt.Errorf("cannot delete member because no member with this ID exists, id=%d", memberID))

		return
	}

	// if correct response send the member back
	c.Header("Content-Type", "application/json")
	// TODO: should this return the deleted member?
	c.JSON(http.StatusOK, nil)
}

// GetAllMembers godoc
// @Summary		Get IDs of all members
// @Description	Get the ID of every member in the database.
// TODO this should eventually be paginated?
// @Tags		members
// @Produce		json
// @Success		200		{array}		models.MemberShortFormDTO
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router		/members	[get]
func (memberController *MemberController) GetAllMembers(c *gin.Context) {
	members, err := memberController.MemberService.GetAllMembers()

	// if there was an error, print it and return status 404: not found
	if err != nil {
		fmt.Println(err)
		utils.ThrowHTTPError(c, http.StatusNotFound, fmt.Errorf("could not retrieve all members: %w", err))

		return
	}

	// if correct response send the member ids and names back
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &members)
}

// GetMemberPosts godoc
// @Summary		Get all posts of this member
// @Description	Get all posts that this member is a collaborator of
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/posts 		[get]
func (memberController *MemberController) GetMemberPosts(_ *gin.Context) {
	// return all the posts
	// that this member is a collaborator/author of
	// TODO implement
}

// GetMemberProjectPosts godoc
// @Summary		Get all project posts of this member
// @Description	Get all project posts that this member is a collaborator of
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/project-posts 		[get]
func (memberController *MemberController) GetMemberProjectPosts(_ *gin.Context) {
	// return all the project posts
	// that this member is a collaborator/author of
	// TODO implement
}

// GetMemberBranches godoc
// @Summary		Get all branches of this member
// @Description	Get all branches that this member is a collaborator of
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/branches 		[get]
func (memberController *MemberController) GetMemberBranches(_ *gin.Context) {
	// return all the branches
	// that this member is a collaborator/author of
	// TODO implement
}

// GetMemberDiscussions godoc
// @Summary		Get all discussions
// @Description	Get all discussions that this member has participated in
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/discussions		[get]
func (memberController *MemberController) GetMemberDiscussions(_ *gin.Context) {
	// returns all the discussions this member is a part of
	// TODO implement
}

// AddMemberSavedPost godoc
// @Summary 	Adds new saved post
// @Description Adds a post to the saved posts of a member
// @Tags 		members
// @Accept  	json
// @Param		memberID		path		string			true	"member ID"
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/saved-posts/{postID} 		[post]
func (memberController *MemberController) AddMemberSavedPost(_ *gin.Context) {

}

// AddMemberSavedProjectPost godoc
// @Summary 	Adds new saved project post
// @Description Adds a project post to the saved project posts of a member
// @Tags 		members
// @Accept  	json
// @Param		memberID		path		string			true	"member ID"
// @Param		postID		path		string			true	"post ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/saved-project-posts/{postID} 		[post]
func (memberController *MemberController) AddMemberSavedProjectPost(_ *gin.Context) {

}

// GetMemberSavedPosts godoc
// @Summary		Get all saved posts of this member
// @Description	Get all posts that this member has saved
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/saved-posts 		[get]
func (memberController *MemberController) GetMemberSavedPosts(_ *gin.Context) {
	// return all saved posts of this member
	// TODO implement
}

// GetMemberProjectPosts godoc
// @Summary		Get all saved project posts of this member
// @Description	Get all project posts that this member has saved
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/members/{memberID}/saved-project-posts 		[get]
func (memberController *MemberController) GetMemberSavedProjectPosts(_ *gin.Context) {
	// return all the project posts that this member has saved
	// TODO implement
}
