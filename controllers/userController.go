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

type UserController struct {
	UserService interfaces.UserService
}

// GetMember godoc
// @Summary 	Get member from database
// @Description Get a member by user ID
// @Accept  	json
// @Param		userID		path		string			true	"user ID"
// @Produce		json
// @Success 	200 		{object}	models.Member
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Router 		/member/{userID}	[get]
func (userController *UserController) GetMember(c *gin.Context) {
	// extract the id of the member
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	// if this caused an error, print it
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid post ID, cannot interpret as integer, id=%s ", userIDStr)})

		return
	}
	// get the user through the service
	member, err := userController.UserService.GetMember(userID)

	// if there was an error, print it and return
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot get member because no user with this ID exists")})

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
// @Success 	200 	{object} 	models.Member
// @Failure		400 	{object} 	utils.HTTPError
// @Router 		/member 		[post]
func (userController *UserController) CreateMember(c *gin.Context) {
	// get the member
	form := forms.MemberCreationForm{}
	// bind the fields of the param to the JSON of the model
	err := c.BindJSON(&form)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind userCreationForm from request body")})

		return
	}

	// create and add to database(not done yet) through the userService
	member := userController.UserService.CreateMember(&form)

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
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Router 		/member 		[put]
func (userController *UserController) UpdateMember(c *gin.Context) {
	// get the new member object
	updatedMember := models.Member{}
	err := c.BindJSON(&updatedMember)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind updated user from request body")})

		return
	}

	// update and add the member to the database
	err = userController.UserService.UpdateMember(&updatedMember)

	// check for errors again
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot update user because no user with this ID exists")})

		return
	}

	// send back a positive response if member updated successfully
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}

// GetCollaborator godoc
// @Summary 	Get collaborator
// @Description Get collaborator by user ID
// @Accept  	json
// @Param		userID		path		string			true	"user ID"
// @Produce		json
// @Success 	200 		{object}	models.Collaborator
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Router 		/collaborator/{userID}	[get]
func (userController *UserController) GetCollaborator(c *gin.Context) {
	// get the user id from the input
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid user ID, cannot interpret as integer, id=%s ", userIDStr)})

		return
	}

	// get the collaborator from the database
	collaborator, err := userController.UserService.GetCollaborator(uint64(userID))

	// check if collaborator found and returned successfully
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot get project user because no user with this ID exists")})

		return
	}

	// if successful, send back the collaborator
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, collaborator)
}

// CreateCollaborator godoc
// @Summary 	Create new collaborator
// @Description Create a new collaborator
// @Accept  	json
// @Param		form	body		forms.CollaboratorCreationForm	true	"Collaborator Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.Collaborator
// @Failure		400 	{object} 	utils.HTTPError
// @Router 		/collaborator 		[post]
func (userController *UserController) CreateCollaborator(c *gin.Context) {
	// extract the fields of the form
	form := forms.CollaboratorCreationForm{}
	err := c.BindJSON(&form)

	// check for errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind CollaboratorCreationForm from request body")})

		return
	}

	// create a collaborator and add to database through the user service
	collaborator := userController.UserService.CreateCollaborator(&form)

	// send back the created collaborator
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &collaborator)
}

// UpdateCollaborator godoc
// @Summary 	Update a collaborator
// @Description Update any number of the fields of a collaborator
// @Accept  	json
// @Param		collaborator	body		models.Collaborator		true	"Updated Collaborator"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Router 		/collaborator 		[put]
func (userController *UserController) UpdateCollaborator(c *gin.Context) {
	// extract the collaborator from the param
	updatedCollaborator := models.PostCollaborator{}
	err := c.BindJSON(&updatedCollaborator)

	// check for errors in the binding
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind updated collaborator from request body")})

		return
	}

	// update collaborator and add to database through the userService
	err = userController.UserService.UpdateCollaborator(&updatedCollaborator)

	// check for errors in database connection
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("cannot update user because no Projectuser with this ID exists")})

		return
	}

	// if updated successfully return an OK response
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)
}
