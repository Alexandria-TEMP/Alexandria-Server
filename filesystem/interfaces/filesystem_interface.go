package filesystem_interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	SetCurrentVersion(versionID, postID uint)
	SaveRepository(c *gin.Context, file *multipart.FileHeader) error
	Unzip() error
	RemoveProjectDirectory() error
	RemoveRepository() error
	GetCurrentDirPath() string
	GetCurrentQuartoDirPath() string
	GetCurrentZipFilePath() string
	GetCurrentRenderDirPath() string
}
