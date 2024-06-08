package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

// @BasePath /api/v2

type BranchController struct {
	BranchService interfaces.BranchService
	RenderService interfaces.RenderService
}

// GetBranch godoc
// @Summary 	Get branch
// @Description Get a branch by branch ID
// @Tags 		branches
// @Accept  	application/json
// @Param		branchID		path		string			true	"Branch ID"
// @Produce		application/json
// @Success 	200 		{object}	models.BranchDTO
// @Failure		400
// @Failure		404
// @Router 		/branches/{branchID}	[get]
func (branchController *BranchController) GetBranch(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%s ", branchIDStr)})

		return
	}

	// get branch and check it exists
	branch, err := branchController.BranchService.GetBranch(uint(branchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, branch.IntoDTO())
}

// CreateBranch godoc
// @Summary 	Create new branch
// @Description Create a new branch linked to a project post.
// @Description Note that Member IDs passed here, get converted to Collaborator IDs.
// @Tags 		branches
// @Accept  	application/json
// @Param		form	body	forms.BranchCreationForm	true	"Branch Creation Form"
// @Produce		application/json
// @Success 	200 	{object} 	models.BranchDTO
// @Failure		404
// @Failure		500
// @Router 		/branches 		[post]
func (branchController *BranchController) CreateBranch(c *gin.Context) {
	// extract branchCreationForm
	form := forms.BranchCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind BranchCreationForm from request body"})

		return
	}

	branch, err404, err500 := branchController.BranchService.CreateBranch(&form)

	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err404.Error()})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err500.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, branch.IntoDTO())
}

// UpdateBranch godoc
// @Summary 	Update branch
// @Description Update any number of the aspects of a branch
// @Tags 		branches
// @Accept  	json
// @Param		branch	body		models.BranchDTO		true	"Updated Branch"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Router 		/branches 		[put]
func (branchController *BranchController) UpdateBranch(c *gin.Context) {
	// extract branchDTO
	dto := models.BranchDTO{}
	err := c.BindJSON(&dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind BranchCreationForm from request body"})

		return
	}
	// update branch
	branch, err := branchController.BranchService.UpdateBranch(dto)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, branch.IntoDTO())
}

// DeleteBranch godoc
// @Summary 	Delete a branch
// @Description Delete a branch with given ID from database
// @Tags 		branches
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/branches/{branchID} 		[delete]
func (branchController *BranchController) DeleteBranch(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%s ", branchIDStr)})

		return
	}

	// delete branch
	if err := branchController.BranchService.DeleteBranch(uint(branchID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	c.Status(http.StatusOK)
}

// GetReviewStatus godoc
// @Summary 	Returns status of all branch reviews
// @Description Returns an array of the statuses of all the reviews of this branch
// @Tags 		branches
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200		{array}		models.ReviewStatus
// @Failure		400
// @Failure		404
// @Router 		/branches/{branchID}/review-statuses	[get]
func (branchController *BranchController) GetReviewStatus(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%s ", branchIDStr)})

		return
	}

	// Get statuses of a branch
	statuses, err := branchController.BranchService.GetReviewStatus(uint(branchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, statuses)
}

// GetReview godoc
// @Summary 	Returns a review
// @Description Returns a review with given ID
// @Tags 		branches
// @Accept  	json
// @Param		reviewID			path		string			true	"review ID"
// @Produce		json
// @Success 	200		{object}	models.ReviewDTO
// @Failure		400
// @Failure		404
// @Router 		/branches/reviews/{reviewID}		[get]
func (branchController *BranchController) GetReview(c *gin.Context) {
	// extract reviewID
	reviewIDStr := c.Param("reviewID")
	reviewID, err := strconv.ParseInt(reviewIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid review ID, cannot interpret as integer, id=%s ", reviewIDStr)})

		return
	}

	// get review
	review, err := branchController.BranchService.GetReview(uint(reviewID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, review)
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
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/branches/{branchID}/reviews		[post]
func (branchController *BranchController) CreateReview(c *gin.Context) {
	// extract ReviewCreationForm
	form := forms.ReviewCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind ReviewCreationForm from request body"})

		return
	}

	// create review and add to branch
	review, err := branchController.BranchService.CreateReview(form)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, review)
}

// UserCanReview godoc
// @Summary 	Returns whether the user is allowed to review this branch
// @Description Returns true if the user fulfills the requirements to review the branch
// @Description Returns false if user is unauthorized to review the branch
// @Tags 		branches
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{object}		boolean
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/branches/{branchID}/can-review/{memberID}		[get]
func (branchController *BranchController) MemberCanReview(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%s ", branchIDStr)})

		return
	}

	// extract memberID
	memberIDStr := c.Param("memberID")
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%s ", memberIDStr)})

		return
	}

	// create review and add to branch
	canReview, err := branchController.BranchService.MemberCanReview(uint(branchID), uint(memberID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, canReview)
}

// GetCollaborator godoc
// @Summary 	Get a branch collaborator by ID
// @Description	Get a branch collaborator by ID, a member who has collaborated on a branch
// @Tags		branches
// @Accept  	json
// @Param		collaboratorID	path	string	true	"Collaborator ID"
// @Produce		json
// @Success 	200 		{object}	models.BranchCollaboratorDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/branches/collaborators/{collaboratorID}	[get]
func (branchController *BranchController) GetBranchCollaborator(_ *gin.Context) {
	// TODO return collaborator by ID
}

// GetRender
// @Summary 	Get the render of a branch
// @Description Get the render of the repository underlying a branch if it exists and has been rendered successfully
// @Tags 		branches
// @Param		branchID	path		string				true	"Branch ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Success		202		{object}	[]byte
// @Failure		400
// @Failure		404
// @Router 		/branches/{branchID}/render	[get]
func (branchController *BranchController) GetRender(c *gin.Context) {
	// extract branchID id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%v ", branchIDStr)})

		return
	}

	// get render filepath
	filePath, err202, err404 := branchController.RenderService.GetRenderFile(uint(branchID))

	// if render is pending return 202 accepted
	if err202 != nil {
		c.String(http.StatusAccepted, "text/plain", []byte("pending"))

		return
	}

	// if render is failed return 404 not found
	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err404.Error()})

		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=render.html")
	c.Header("Content-Type", "text/html")
	c.File(filePath)
}

// GetProject godoc specs are subject to change
// @Summary 	Get the repository of a branch
// @Description Get the entire zipped repository of a branch
// @Tags 		branches
// @Param		branchID	path		string				true	"Branch ID"
// @Produce		application/zip
// @Success 	200		{object}	[]byte
// @Failure		400
// @Failure		404
// @Router 		/branches/{branchID}/repository	[get]
func (branchController *BranchController) GetProject(c *gin.Context) {
	// extract branch id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%v ", branchIDStr)})

		return
	}

	// get repository filepath
	filePath, err := branchController.BranchService.GetProject(uint(branchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=quarto_project.zip")
	c.Header("Content-Type", "application/zip")
	c.File(filePath)
}

// UploadProject
// @Summary 	Upload a new project version to a branch
// @Description Upload a new project version to a specific, preexisting, branch as a zipped quarto project
// @Tags 		branches
// @Accept  	multipart/form-data
// @Param		branchID		path		string			true	"Branch ID"
// @Param		file			formData	file			true	"Repository to create"
// @Produce		application/json
// @Success 	200
// @Failure		400
// @Failure		500
// @Router 		/branches/{branchID}		[post]
func (branchController *BranchController) UploadProject(c *gin.Context) {
	// extract file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file found"})

		return
	}

	// extract branch id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%v ", branchIDStr)})

		return
	}

	// Create commit on branch with new files
	err = branchController.BranchService.UploadProject(c, file, uint(branchID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// response
	c.Status(http.StatusOK)
}

// GetFiletree godoc specs are subject to change
// @Summary 	Get the filetree of a project
// @Description Get the filetree of a project of a branch
// @Tags 		branches
// @Param		branchID	path		string				true	"Branch ID"
// @Produce		application/json
// @Success 	200		{object}	map[string]int64
// @Failure		400
// @Failure		404
// @Failure		500
// @Router 		/branches/{branchID}/tree		[get]
func (branchController *BranchController) GetFiletree(c *gin.Context) {
	// extract branchID id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%v ", branchIDStr)})

		return
	}

	fileTree, err404, err500 := branchController.BranchService.GetFiletree(uint(branchID))

	if err404 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err404.Error()})

		return
	}

	if err500 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err500.Error()})

		return
	}

	// response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, fileTree)
}

// GetFileFromProject godoc specs are subject to change
// @Summary 	Get a file from a project
// @Description Get the contents of a single file from a project of a branch
// @Tags 		branches
// @Param		branchID	path		string				true	"Branch ID"
// @Param		filepath	path		string				true	"Filepath"
// @Produce		application/octet-stream
// @Success 	200		{object}	[]byte
// @Failure		404
// @Failure		500
// @Router 		/branches/{branchID}/file/{filepath}	[get]
func (branchController *BranchController) GetFileFromProject(c *gin.Context) {
	// extract branchID id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID, cannot interpret as integer, id=%v ", branchIDStr)})

		return
	}

	relFilepath := c.Param("filepath")
	absFilepath, err := branchController.BranchService.GetFileFromProject(uint(branchID), relFilepath)

	// if files doesnt exist return 404 not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// get the file info
	fileContentType, err1 := mimetype.DetectFile(absFilepath)
	fileData, err2 := os.Open(absFilepath)
	fileInfo, err3 := fileData.Stat()

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})

		return
	}

	defer fileData.Close()

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	c.Header("Content-Type", fileContentType.String())
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(absFilepath)
}

// GetDiscussions godoc
// @Summary Returns all level 1 discussions associated with the version
// @Description Returns all discussions on this version that are not a reply to another discussion
// @Description Endpoint is offset-paginated
// @Tags 		branches
// @Param		branchID	path		string			true	"Branch ID"
// @Param 		page		query		uint			false	"page query"
// @Param		pageSize	query		uint			false	"page size"
// @Produce		application/json
// @Success 	200		{array}		models.DiscussionDTO
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/branches/{branchID}/discussions 	[get]
func (branchController *BranchController) GetDiscussions(_ *gin.Context) {

}
