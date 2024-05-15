package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -source=./versionService_interface.go -destination=../../mocks/versionService_mock.go

type VersionService interface {
	SaveRepository(c *gin.Context, file *multipart.FileHeader, versionID uint, postID uint) error
	CreateVersion() *models.Version
	GetFilesystem() filesystem.Filesystem
}
