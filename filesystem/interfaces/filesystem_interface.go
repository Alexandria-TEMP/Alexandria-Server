package filesysteminterface

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	SetCurrentVersion(versionID uint)
	SaveRepository(c *gin.Context, file *multipart.FileHeader) error
	Unzip() error
	RenderExists() (bool, string)
	RemoveRepository() error
	GetCurrentDirPath() string
	GetCurrentQuartoDirPath() string
	GetCurrentZipFilePath() string
	GetCurrentRenderDirPath() string
	GetFileTree() (map[string]int64, error)
}
