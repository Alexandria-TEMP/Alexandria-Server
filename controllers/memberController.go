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
	c.JSON(http.StatusOK, member.IntoDTO())
}

// CreateMember godoc
// @Summary 	Create a new member
// @Description Create a new member from the given fields.
// @Description The member must have a unique email address, which isn't associated with any other accounts.
// @Description They are automatically logged in, and an access + refresh token pair is returned alongside the member
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

	// if there is an error, return a 404 not found status
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot bind tag ids from request body: %s", err)})

		return
	}

	tagContainer := models.ScientificFieldTagContainer{
		ScientificFieldTags: tagArray,
	}

	// create and add to database through the memberService
	loggedInMember, err := memberController.MemberService.CreateMember(&form, &tagContainer)

	// if the member service throws an error, return a 400 Bad request status
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create member: %s", err)})

		return
	}

	// send back a positive response 200 status with the created member
	c.JSON(http.StatusOK, loggedInMember)
}

// UpdateMember godoc
// @Summary 	Update a member
// @Description Update the fields of a member
// @Tags 		members
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		member	body		models.MemberDTO		true	"Updated member"
// @Produce		json
// @Success 	200
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/members 		[put]
func (memberController *MemberController) UpdateMember(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// DeleteMember godoc
// @Summary 	Delete a member
// @Description Delete a member with given ID from database
// @Tags 		members
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
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
	memberShortDTOs, err := memberController.MemberService.GetAllMembers()

	// if there was an error, print it and return status 404: not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not retrieve all members: %s", err)})

		return
	}

	// if correct response send the member ids and names back
	c.JSON(http.StatusOK, &memberShortDTOs)
}

// LoginMember godoc
// @Summary		Logs a member in
// @Description	Logs a member in based on email and password and returns an access and refresh token.
// @Tags 		members
// @Accept 		json
// @Param		member	body		forms.MemberAuthForm		true	"Member Authentication Form"
// @Produce		json
// @Success 	200		{object}	models.LoggedInMemberDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/login 		[post]
func (memberController *MemberController) LoginMember(c *gin.Context) {
	form := &forms.MemberAuthForm{}
	// bind the fields of the param to the JSON of the model
	err := c.BindJSON(form)

	// if there is an error, return a 400 bad request status
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind MemberAuthForm from request body"})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form fails validation"})
		return
	}

	loggedInMember, err := memberController.MemberService.LogInMember(form)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to log in member: %s", err.Error())})

		return
	}

	c.JSON(http.StatusOK, loggedInMember)
}

// RefreshToken godoc
// @Summary		Refreshes the access token.
// @Description	Refreshes the access token with a refresh token.
// @Tags 		members
// @Accept 		json
// @Param		member	body		forms.TokenRefreshForm		true	"Token Refresh Form"
// @Produce		json
// @Success 	200		{object}	models.TokenPairDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/members/token	[post]
func (memberController *MemberController) RefreshToken(c *gin.Context) {
	form := &forms.TokenRefreshForm{}
	// bind the fields of the param to the JSON of the model
	err := c.BindJSON(form)

	// if there is an error, return a 400 bad request status
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind TokenRefreshForm from request body"})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form fails validation"})
		return
	}

	tokenPair, err := memberController.MemberService.RefreshToken(form)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to refresh access token: %s", err.Error())})

		return
	}

	c.JSON(http.StatusOK, tokenPair)
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
// @Param 		Authorization header string true "Access Token"
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
// @Param 		Authorization header string true "Access Token"
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
