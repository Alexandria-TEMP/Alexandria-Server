package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type BranchController struct {
}

// GetBranch godoc
// @Summary 	Get branch
// @Description Get a branch by branch ID
// @Tags 		branches
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
// @Description Create a new branch linked to a project post.
// @Description Note that Member IDs passed here, get converted to Collaborator IDs.
// @Tags 		branches
// @Accept  	json
// @Param		form	body	forms.BranchCreationForm	true	"Branch Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.BranchDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/branches 		[post]
func (branchController *BranchController) CreateBranch(_ *gin.Context) {

}

// UpdateBranch godoc
// @Summary 	Update branch
// @Description Update any number of the aspects of a branch
// @Tags 		branches
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
// @Tags 		branches
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
// @Tags 		branches
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200		{array}		models.BranchReviewStatus
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID}/review-statuses		[get]
func (branchController *BranchController) GetReviewStatus(_ *gin.Context) {
	// delete method goes here
}

// GetReview godoc
// @Summary 	Returns a branch review by ID
// @Description Returns a review of a branch with the given ID
// @Tags 		branches
// @Accept  	json
// @Param		reviewID			path		string			true	"review ID"
// @Produce		json
// @Success 	200		{object}	models.BranchReviewDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/reviews/{reviewID}		[get]
func (branchController *BranchController) GetReview(_ *gin.Context) {

}

// CreateReview godoc
// @Summary 	Adds a review to a branch
// @Description Adds a review to a branch
// @Tags 		branches
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
// @Tags 		branches
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		memberID			path		string			true	"user ID"
// @Produce		json
// @Success 	200		{object}		boolean
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/branches/{branchID}/can-review/{memberID}		[get]
func (branchController *BranchController) UserCanReview(_ *gin.Context) {

}

// GetCollaborator godoc
// @Summary 	Get a branch collaborator by ID
// @Description	Get a branch collaborator by ID, a member who has collaborated on a branch
// @Tags		branches
// @Accept  	json
// @Param		collaboratorID	path	string	true	"Collaborator ID"
// @Produce		json
// @Success 	200 		{object}	models.BranchCollaboratorDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/branches/collaborators/{collaboratorID}	[get]
func (branchController *BranchController) GetBranchCollaborator(_ *gin.Context) {
	// TODO return collaborator by ID
}
