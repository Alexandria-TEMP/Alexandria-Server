package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type BranchController struct {
	BranchService interfaces.BranchService
}

// GetBranch godoc
// @Summary 	Get branch
// @Description Get a branch by branch ID
// @Accept  	json
// @Param		branchID		path		string			true	"Branch ID"
// @Produce		json
// @Success 	200 		{object}	models.BranchDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}	[get]
func (branchController *BranchController) GetBranch(_ *gin.Context) {

}

// CreateBranch godoc
// @Summary 	Create new branch
// @Description Create a new question or discussion branch
// @Accept  	json
// @Param		form	body	forms.BranchCreationForm	true	"Branch Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.BranchDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/branches 		[post]
func (branchController *BranchController) CreateBranch(c *gin.Context) {
	// extract branchCreationForm
	form := forms.BranchCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind BranchCreationForm from request body"})

		return
	}

	branch, err404, err500 := branchController.BranchService.CreateBranch(form)

	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, &branch)
}

// UpdateBranch godoc
// @Summary 	Update branch
// @Description Update any number of the aspects of a branch
// @Accept  	json
// @Param		branch	body		models.BranchDTO		true	"Updated Branch"
// @Produce		json
// @Success 	200
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/branches 		[put]
func (branchController *BranchController) UpdateBranch(_ *gin.Context) {

}

// DeleteBranch godoc
// @Summary 	Delete a branch
// @Description Delete a branch with given ID from database
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID} 		[delete]
func (branchController *BranchController) DeleteBranch(_ *gin.Context) {
	// delete method goes here
}

// GetReviewStatus godoc
// @Summary 	Returns status of all branch reviews
// @Description Returns an array of the statuses of all the reviews of this branch
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200		{array}		string
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID}/reviews		[get]
func (branchController *BranchController) GetReviewStatus(_ *gin.Context) {
	// delete method goes here
}

// GetReview godoc
// @Summary 	Returns a review of a branch
// @Description Returns a review with the given ID of the branch with the given ID
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		reviewID			path		string			true	"review ID"
// @Produce		json
// @Success 	200		{object}	models.ReviewDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID}/reviews/{reviewID}		[get]
func (branchController *BranchController) GetReview(_ *gin.Context) {

}

// CreateReview godoc
// @Summary 	Adds a review to a branch
// @Description Adds a review to a branch
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		form	body	forms.ReviewCreationForm	true	"review creation form"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID}/reviews		[post]
func (branchController *BranchController) CreateReview(_ *gin.Context) {

}

// UserCanReview godoc
// @Summary 	Returns whether the user is allowed to review this branch
// @Description Returns true if the user fulfills the requirements to review the branch
// @Description Returns false if user is unauthorized to review the branch
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		userID			path		string			true	"user ID"
// @Produce		json
// @Success 	200		{array}		boolean
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID}/can-review/{userID}		[get]
func (branchController *BranchController) UserCanReview(_ *gin.Context) {

}

// - `/branches`
//   - `POST`
//   - `PUT` (?)
//   - `/:id` `GET`
//   - `/:id` `DELETE`
//   - `/:id/reviews` `GET` (gets acceptance status of all reviews)
//   - `/:id/reviews/:id` `GET` (gets specific review)
//   - `/:id/reviews` `POST` (does the merge - make sure to refresh page)
//   - `/:id/reviews/can-review` `GET` (utility endpoint for front-end - allowed to review?)
