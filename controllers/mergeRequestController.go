package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type MergeRequestController struct {
}

// GetMergeRequest godoc
// @Summary 	Get merge request
// @Description Get a merge request by merge request ID
// @Tags 		merge-requests
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"MergeRequest ID"
// @Produce		json
// @Success 	200 		{object}	models.MergeRequestDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}	[get]
func (mergeRequestController *MergeRequestController) GetMergeRequest(_ *gin.Context) {

}

// CreateMergeRequest godoc
// @Summary 	Create new merge request
// @Description Create a new merge request linked to a project post.
// @Description Note that Member IDs passed here, get converted to Collaborator IDs.
// @Tags 		merge-requests
// @Accept  	json
// @Param		form	body	forms.MergeRequestCreationForm	true	"MergeRequest Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.MergeRequestDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/merge-requests 		[post]
func (mergeRequestController *MergeRequestController) CreateMergeRequest(_ *gin.Context) {

}

// UpdateMergeRequest godoc
// @Summary 	Update merge request
// @Description Update any number of the aspects of a merge request
// @Tags 		merge-requests
// @Accept  	json
// @Param		merge request	body		models.MergeRequestDTO		true	"Updated MergeRequest"
// @Produce		json
// @Success 	200
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/merge-requests 		[put]
func (mergeRequestController *MergeRequestController) UpdateMergeRequest(_ *gin.Context) {

}

// DeleteMergeRequest godoc
// @Summary 	Delete a merge request
// @Description Delete a merge request with given ID from database
// @Tags 		merge-requests
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID} 		[delete]
func (mergeRequestController *MergeRequestController) DeleteMergeRequest(_ *gin.Context) {
	// delete method goes here
}

// GetReviewStatus godoc
// @Summary 	Returns status of all merge request reviews
// @Description Returns an array of the statuses of all the reviews of this merge request
// @Tags 		merge-requests
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Produce		json
// @Success 	200		{array}		string
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/reviews		[get]
func (mergeRequestController *MergeRequestController) GetReviewStatus(_ *gin.Context) {
	// delete method goes here
}

// GetReview godoc
// @Summary 	Returns a merge request review by ID
// @Description Returns a review of a merge request with the given ID
// @Tags 		merge-requests
// @Accept  	json
// @Param		reviewID			path		string			true	"review ID"
// @Produce		json
// @Success 	200		{object}	models.ReviewDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/reviews/{reviewID}		[get]
func (mergeRequestController *MergeRequestController) GetReview(_ *gin.Context) {

}

// CreateReview godoc
// @Summary 	Adds a review to a merge request
// @Description Adds a review to a merge request
// @Tags 		merge-requests
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Param		form	body	forms.ReviewCreationForm	true	"review creation form"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/reviews		[post]
func (mergeRequestController *MergeRequestController) CreateReview(_ *gin.Context) {

}

// UserCanReview godoc
// @Summary 	Returns whether the user is allowed to review this merge request
// @Description Returns true if the user fulfills the requirements to review the merge request
// @Description Returns false if user is unauthorized to review the merge request
// @Tags 		merge-requests
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Param		userID			path		string			true	"user ID"
// @Produce		json
// @Success 	200		{array}		boolean
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/can-review/{userID}		[get]
func (mergeRequestController *MergeRequestController) UserCanReview(_ *gin.Context) {

}

// - `/merge-requests`
//   - `POST`
//   - `PUT` (?)
//   - `/:id` `GET`
//   - `/:id` `DELETE`
//   - `/:id/reviews` `GET` (gets acceptance status of all reviews)
//   - `/:id/reviews/:id` `GET` (gets specific review)
//   - `/:id/reviews` `POST` (does the merge - make sure to refresh page)
//   - `/:id/reviews/can-review` `GET` (utility endpoint for front-end - allowed to review?)

// GetCollaborator godoc
// @Summary 	Get a merge request collaborator by ID
// @Description	Get a merge request collaborator by ID, a member who has collaborated on a merge request
// @Tags		merge-requests
// @Accept  	json
// @Param		collaboratorID	path	string	true	"Collaborator ID"
// @Produce		json
// @Success 	200 		{object}	models.MergeRequestCollaboratorDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/merge-requests/collaborators/{collaboratorID}	[get]
func (mergeRequestController *MergeRequestController) GetMergeRequestCollaborator(_ *gin.Context) {
	// TODO return collaborator by ID
}
