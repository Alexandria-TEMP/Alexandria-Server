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
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Failure		500			{object} 	utils.HTTPError
// @Router 		/members/{memberID}	[get]
func (memberController *MemberController) GetMember(c *gin.Context) {
	// extract the id of the member
	memberIDStr := c.Param("memberID")
	initmemberID, err := strconv.ParseUint(memberIDStr, 10, 64)
	// if this caused an error, print it and return status 400: bad input
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid member ID, cannot interpret '%s' as integer: %s", memberIDStr, err)})

		return
	}

	// cast member ID as uint instead of uint64, because database only accepts those
	memberID := uint(initmemberID)

	// get the member through the service
	member, err := memberController.MemberService.GetMember(memberID)

	// if there was an error, print it and return status 404: not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not get member with ID %d: %s", memberID, err)})

		return
	}

	// if correct response send the member back
	c.JSON(http.StatusOK, member)
}

// CreateMember godoc
// @Summary 	Create a new member
// @Description Create a new member from the given fields
// @Tags 		members
// @Accept  	json
// @Param		form	body	forms.MemberCreationForm	true	"Member Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.MemberDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members 		[post]
func (memberController *MemberController) CreateMember(c *gin.Context) {
	form := forms.MemberCreationForm{}
	// bind the fields of the param to the JSON of the model
	err := c.BindJSON(&form)

	// if there is an error, return a 400 bad request status
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind memberCreationForm from request body: %s", err)})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	// get array of strings, create array of tags
	tagIDs := form.ScientificFieldTagIDs

	// getting the tags from tag service
	tagArray, err := memberController.TagService.GetTagsFromIDs(tagIDs)

	tagContainer := models.ScientificFieldTagContainer{
		ScientificFieldTags: tagArray,
	}

	// if there is an error, return a 404 not found status
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot bind tag ids from request body: %s", err)})

		return
	}

	// create and add to database through the memberService
	member, err := memberController.MemberService.CreateMember(&form, &tagContainer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create member: %s", err)})

		return
	}

	// send back a positive response 200 status with the created member
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, member)
}

// UpdateMember godoc
// @Summary 	Update a member
// @Description Update the fields of a member
// @Tags 		members
// @Accept  	json
// @Param		member	body		models.MemberDTO		true	"Updated member"
// @Produce		json
// @Success 	200
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members 		[put]
func (memberController *MemberController) UpdateMember(c *gin.Context) {
	// get the new member object
	updatedMember := models.MemberDTO{}
	err := c.BindJSON(&updatedMember)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind updated member from request body: %s", err)})

		return
	}

	// get array of strings, create array of tags
	tagIDs := updatedMember.ScientificFieldTagIDs
	// call the method from the tag service
	tagArray, err := memberController.TagService.GetTagsFromIDs(tagIDs)
	tagContainer := models.ScientificFieldTagContainer{
		ScientificFieldTags: tagArray,
	}

	// if there is an error, return a 400 bad request status
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind tag IDs from request body: %s", err)})

		return
	}

	err = memberController.MemberService.UpdateMember(&updatedMember, &tagContainer)
	// check for errors again
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to update member: %s", err)})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID} 		[delete]
func (memberController *MemberController) DeleteMember(c *gin.Context) {
	// extract the id of the member
	memberIDStr := c.Param("memberID")
	initmemberID, err := strconv.ParseUint(memberIDStr, 10, 64)

	// if this caused an error, print it and return status 400: bad input
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid member ID, cannot interpret '%s' as integer: %s ", memberIDStr, err)})

		return
	}

	// cast member ID as uint instead of uint64, because database only accepts those
	memberID := uint(initmemberID)

	// get the member through the service
	err = memberController.MemberService.DeleteMember(memberID)

	// if there was an error, print it and return status 404: not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot delete member because no member with ID '%d' exists: %s", memberID, err)})

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
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/members	[get]
func (memberController *MemberController) GetAllMembers(c *gin.Context) {
	members, err := memberController.MemberService.GetAllMembers()

	// if there was an error, print it and return status 404: not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not retrieve all members: %s", err)})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/posts 		[get]
func (memberController *MemberController) GetMemberPosts(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetMemberProjectPosts godoc
// @Summary		Get all project posts of this member
// @Description	Get all project posts that this member is a collaborator of
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/project-posts 		[get]
func (memberController *MemberController) GetMemberProjectPosts(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetMemberBranches godoc
// @Summary		Get all branches of this member
// @Description	Get all branches that this member is a collaborator of
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/branches 		[get]
func (memberController *MemberController) GetMemberBranches(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetMemberDiscussions godoc
// @Summary		Get all discussions
// @Description	Get all discussions that this member has participated in
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/discussions		[get]
func (memberController *MemberController) GetMemberDiscussions(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/saved-posts/{postID} 		[post]
func (memberController *MemberController) AddMemberSavedPost(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/saved-project-posts/{postID} 		[post]
func (memberController *MemberController) AddMemberSavedProjectPost(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetMemberSavedPosts godoc
// @Summary		Get all saved posts of this member
// @Description	Get all posts that this member has saved
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/saved-posts 		[get]
func (memberController *MemberController) GetMemberSavedPosts(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetMemberProjectPosts godoc
// @Summary		Get all saved project posts of this member
// @Description	Get all project posts that this member has saved
// @Tags 		members
// @Accept 		json
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{array}		uint
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members/{memberID}/saved-project-posts 		[get]
func (memberController *MemberController) GetMemberSavedProjectPosts(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
