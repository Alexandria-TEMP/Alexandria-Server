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
	BranchService             interfaces.BranchService
	RenderService             interfaces.RenderService
	BranchCollaboratorService interfaces.BranchCollaboratorService
}

// GetBranch godoc
// @Summary 	Get branch
// @Description Get a branch by branch ID
// @Tags 		branches
// @Accept  	application/json
// @Param		branchID		path		string			true	"Branch ID"
// @Produce		application/json
// @Success 	200 		{object}	models.BranchDTO
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Router 		/branches/{branchID}	[get]
func (branchController *BranchController) GetBranch(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

		return
	}

	// get branch and check it exists
	branch, err := branchController.BranchService.GetBranch(uint(branchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, branch.IntoDTO())
}

// CreateBranch godoc
// @Summary 	Create new branch
// @Description Create a new branch linked to a project post.
// @Description Note that Member IDs passed here, get converted to Collaborator IDs.
// @Tags 		branches
// @Accept  	application/json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.BranchCreationForm	true	"Branch Creation Form"
// @Produce		application/json
// @Success 	200 	{object} 	models.BranchDTO
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches 		[post]
func (branchController *BranchController) CreateBranch(c *gin.Context) {
	// extract branchCreationForm
	form := forms.BranchCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind BranchCreationForm from request body: %s", err)})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

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
	c.JSON(http.StatusOK, branch.IntoDTO())
}

// DeleteBranch godoc
// @Summary 	Delete a branch
// @Description Delete a branch with given ID from database
// @Tags 		branches
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches/{branchID} 		[delete]
func (branchController *BranchController) DeleteBranch(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

		return
	}

	// delete branch
	if err := branchController.BranchService.DeleteBranch(uint(branchID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	c.Status(http.StatusOK)
}

// GetAllBranchReviewStatuses godoc
// @Summary 	Returns status of all branch reviews
// @Description Returns an array of the statuses of all the reviews of this branch
// @Tags 		branches
// @Accept  	json
// @Param		branchID		path		string			true	"branch ID"
// @Produce		json
// @Success 	200		{array}		models.BranchOverallReviewStatus
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/review-statuses	[get]
func (branchController *BranchController) GetAllBranchReviewStatuses(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

		return
	}

	// Get statuses of a branch
	statuses, err := branchController.BranchService.GetAllBranchReviewStatuses(uint(branchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, statuses)
}

// GetReview godoc
// @Summary 	Returns a branch review
// @Description Returns a branch review with given ID
// @Tags 		branches
// @Accept  	json
// @Param		reviewID			path		string			true	"branchreview ID"
// @Produce		json
// @Success 	200		{object}	models.BranchReviewDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/branches/reviews/{reviewID}		[get]
func (branchController *BranchController) GetReview(c *gin.Context) {
	// extract reviewID
	reviewIDStr := c.Param("reviewID")
	reviewID, err := strconv.ParseInt(reviewIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branchreview ID '%s', cannot interpret as integer: %s", reviewIDStr, err)})

		return
	}

	// get branchreview
	branchreview, err := branchController.BranchService.GetReview(uint(reviewID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, branchreview.IntoDTO())
}

// CreateReview godoc
// @Summary 	Adds a branchreview to a branch
// @Description Adds a branchreview to a branch
// @Tags 		branches
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		form	body	forms.ReviewCreationForm	true	"branchreview creation form"
// @Produce		json
// @Success 	200		{object}	models.BranchReviewDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches/reviews		[post]
func (branchController *BranchController) CreateReview(c *gin.Context) {
	// extract ReviewCreationForm
	form := forms.ReviewCreationForm{}
	err := c.BindJSON(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("cannot bind ReviewCreationForm from request body: %s", err)})

		return
	}

	if !form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate form"})

		return
	}

	// create branchreview and add to branch
	branchreview, err := branchController.BranchService.CreateReview(form)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, branchreview.IntoDTO())
}

// UserCanReview godoc
// @Summary 	Returns whether the user is allowed to branchreview this branch
// @Description Returns true if the user fulfills the requirements to branchreview the branch
// @Description Returns false if user is unauthorized to branchreview the branch
// @Tags 		branches
// @Accept  	json
// @Param 		Authorization header string true "Access Token"
// @Param		branchID		path		string			true	"branch ID"
// @Param		memberID		path		string			true	"member ID"
// @Produce		json
// @Success 	200		{object}		boolean
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/can-review/{memberID}		[get]
func (branchController *BranchController) MemberCanReview(c *gin.Context) {
	// extract branchID
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

		return
	}

	// extract memberID
	memberIDStr := c.Param("memberID")
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", memberIDStr, err)})

		return
	}

	// create branchreview and add to branch
	canReview, err := branchController.BranchService.MemberCanReview(uint(branchID), uint(memberID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
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
// @Failure		400			{object} 	utils.HTTPError
// @Failure		404			{object} 	utils.HTTPError
// @Router 		/branches/collaborators/{collaboratorID}	[get]
func (branchController *BranchController) GetBranchCollaborator(c *gin.Context) {
	// extract collaboratorID id
	collaboratorIDStr := c.Param("collaboratorID")
	collaboratorID, err := strconv.ParseUint(collaboratorIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", collaboratorIDStr, err)})

		return
	}

	collaborator, err := branchController.BranchCollaboratorService.GetBranchCollaborator(uint(collaboratorID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, collaborator.IntoDTO())
}

// GetRender
// @Summary 	Get the render of a branch
// @Description Get the render of the repository underlying a branch if it exists and has been rendered successfully
// @Tags 		branches
// @Param		branchID	path		string				true	"Branch ID"
// @Produce		text/html
// @Success 	200		{object}	[]byte
// @Success		202		{object}	[]byte
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/render	[get]
func (branchController *BranchController) GetRender(c *gin.Context) {
	// extract branchID id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/repository	[get]
func (branchController *BranchController) GetProject(c *gin.Context) {
	// extract branch id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

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
// @Description Upload a new project version to a specific, preexisting, branch as a zipped quarto project.
// @Description Specifically, this zip should contain all of the contents of the project at its root, not in a subdirectory.
// @Description Call this after you create a post, and supply it with the actual post contents.
// @Tags 		branches
// @Accept  	multipart/form-data
// @Param 		Authorization header string true "Access Token"
// @Param		branchID		path		string			true	"Branch ID"
// @Param		file			formData	file			true	"Repository to create"
// @Produce		application/json
// @Success 	200
// @Failure		400		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/upload		[post]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/tree		[get]
func (branchController *BranchController) GetFiletree(c *gin.Context) {
	// extract branchID id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

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
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router 		/branches/{branchID}/file/{filepath}	[get]
func (branchController *BranchController) GetFileFromProject(c *gin.Context) {
	// extract branchID id
	branchIDStr := c.Param("branchID")
	branchID, err := strconv.ParseUint(branchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid branch ID '%s', cannot interpret as integer: %s", branchIDStr, err)})

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

	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read file: %s", err1)})

		return
	}

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read file: %s", err2)})

		return
	}

	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read file: %s", err3)})

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
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Failure		500		{object} 	utils.HTTPError
// @Router		/branches/{branchID}/discussions 	[get]
func (branchController *BranchController) GetDiscussions(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// GetClosedBranch godoc
// @Summary Returns a closed branch
// @Description Returns a closed branch given an id
// @Tags 		branches
// @Param		closedBranchID	path		string			true	"Closed Branch ID"
// @Produce		application/json
// @Success 	200		{object}	models.ClosedBranchDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router		/branches/closed/{closedBranchID}		[get]
func (branchController *BranchController) GetClosedBranch(c *gin.Context) {
	// extract branchID id
	closedBranchIDStr := c.Param("closedBranchID")
	closedBranchID, err := strconv.ParseUint(closedBranchIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid closed branch ID '%s', cannot interpret as integer: %s", closedBranchIDStr, err)})

		return
	}

	closedBranch, err := branchController.BranchService.GetClosedBranch(uint(closedBranchID))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	// response
	c.JSON(http.StatusOK, closedBranch.IntoDTO())
}

// GetAllBranchCollaborators godoc
// @Summary 	Get all branch collaborators of a branch
// @Description Returns all branch collaborators of the branch with the given ID
// @Tags 		branches
// @Param		branchID	path		string			true	"Branch ID"
// @Produce		application/json
// @Success 	200		{array}		models.BranchCollaboratorDTO
// @Failure		400		{object} 	utils.HTTPError
// @Failure		404		{object} 	utils.HTTPError
// @Router		/branches/collaborators/all/{branchID}		[get]
func (branchController *BranchController) GetAllBranchCollaborators(c *gin.Context) {
	// Get branch ID from path param
	branchIDString := c.Param("branchID")

	branchID, err := strconv.ParseUint(branchIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse branch ID '%s' as unsigned integer: %s", branchIDString, err)})

		return
	}

	// Get the branch itself
	branch, err := branchController.BranchService.GetBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get branch with ID %d: %s", branchID, err)})

		return
	}

	branchCollaborators := branch.Collaborators

	// Turn each branch collaborator into a DTO
	branchCollaboratorDTOs := make([]*models.BranchCollaboratorDTO, len(branchCollaborators))

	for i, branchCollaborator := range branchCollaborators {
		branchCollaboratorDTO := branchCollaborator.IntoDTO()
		branchCollaboratorDTOs[i] = &branchCollaboratorDTO
	}

	c.JSON(http.StatusOK, branchCollaboratorDTOs)
}
