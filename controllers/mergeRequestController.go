package controllers

import "github.com/gin-gonic/gin"

// @BasePath /api/v2

type MergeRequestController struct {
}

// GetMergeRequest godoc
// @Summary 	Get merge request
// @Description Get a merge request by merge request ID
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"MergeRequest ID"
// @Produce		json
// @Success 	200 		{object}	models.MergeRequestDTO
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}	[get]
func (mergeRequestController *MergeRequestController) GetMergeRequest(c *gin.Context) {

}

// CreateMergeRequest godoc
// @Summary 	Create new merge request
// @Description Create a new question or discussion merge request
// @Accept  	json
// @Param		form	body	forms.MergeRequestCreationForm	true	"MergeRequest Creation Form"
// @Produce		json
// @Success 	200 	{object} 	models.MergeRequestDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		500 	{object} 	utils.HTTPError
// @Router 		/merge-requests 		[post]
func (mergeRequestController *MergeRequestController) CreateMergeRequest(c *gin.Context) {

}

// UpdateMergeRequest godoc
// @Summary 	Update merge request
// @Description Update any number of the aspects of a merge request
// @Accept  	json
// @Param		merge request	body		models.MergeRequestDTO		true	"Updated MergeRequest"
// @Produce		json
// @Success 	200
// @Failure		400 		{object} 	utils.HTTPError
// @Failure		404 		{object} 	utils.HTTPError
// @Failure		500 		{object} 	utils.HTTPError
// @Router 		/merge-requests 		[put]
func (mergeRequestController *MergeRequestController) UpdateMergeRequest(c *gin.Context) {

}

// DeleteMergeRequest godoc
// @Summary 	Delete a merge request
// @Description Delete a merge request with given ID from database
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Produce		json
// @Success 	200
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID} 		[delete]
func (mergeRequestController *MergeRequestController) DeleteMergeRequest(c *gin.Context) {
	//delete method goes here
}


// GetReviewStatus godoc
// @Summary 	Returns status of all merge request reviews
// @Description Returns an array of the statuses of all the reviews of this merge request
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Produce		json
// @Success 	200		{array}		string
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/reviews		[get]
func (mergeRequestController *MergeRequestController) GetReviewStatus(c *gin.Context) {
	//delete method goes here
}


// GetReview godoc
// @Summary 	Returns a review of a merge request
// @Description Returns a review with the given ID of the merge request with the given ID
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Param		reviewID			path		string			true	"review ID"
// @Produce		json
// @Success 	200		{object}	models.ReviewDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/reviews/{reviewID}		[get]
func (mergeRequestController *MergeRequestController) GetReview(c *gin.Context) {
	
}


// CreateReview godoc
// @Summary 	Adds a review to a merge request
// @Description Adds a review to a merge request
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Produce		json
// @Success 	200		{object}	models.ReviewDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/reviews		[post]
func (mergeRequestController *MergeRequestController) CreateReview(c *gin.Context) {
	
}

// UserCanReview godoc
// @Summary 	Returns whether the user is allowed to review this merge request
// @Description Returns true if the user fulfills the requirements to review the merge request
// @Description Returns false if user is unauthorized to review the merge request
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Param		userID			path		string			true	"user ID"
// @Produce		json
// @Success 	200		{array}		boolean
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/can-review/{userID}		[get]
func (mergeRequestController *MergeRequestController) UserCanReview(c *gin.Context) {
	
}

// MergeMergeRequest godoc
// @Summary 	Merges the merge request into parent post
// @Description Merges the merge request with the given id into the respective project post
// @Accept  	json
// @Param		mergeRequestID		path		string			true	"merge request ID"
// @Produce		json
// @Success 	200		{object}	models.ProjectPostDTO
// @Failure		400 	{object} 	utils.HTTPError
// @Failure		404 	{object} 	utils.HTTPError
// @Failure		500		{object}	utils.HTTPError
// @Router 		/merge-requests/{mergeRequestID}/merge		[put]
func (mergeRequestController *MergeRequestController) MergeMergeRequest(c *gin.Context) {
	
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