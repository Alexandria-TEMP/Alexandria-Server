package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./versionService_interface.go -destination=../../mocks/versionService_mock.go

type VersionService interface {
	// CreateVersion orchestrates the version creation.
	// 1. creates a new version with the pending render status.
	// 2. saves the file to its directory
	// 3. return a status 200 to client
	// 4. unzip file
	// 5. render project and update render status
	// TODO: persist data
	CreateVersion(c *gin.Context, file *multipart.FileHeader, postID uint) (*models.Version, error)

	// GetRender returns filepath of rendered repository.
	// Error 1 is for status 202.
	// Error 2 is for status 404.
	GetRenderFile(versionID, postID uint) (string, error, error)

	// GetRender returns filepath of zipped repository
	GetRepositoryFile(versionID, postID uint) (string, error)
}
