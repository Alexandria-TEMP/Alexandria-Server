package filesysteminterface

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	SetCurrentVersion(versionID, postID uint)
	SaveRepository(c *gin.Context, file *multipart.FileHeader) error
	Unzip() error
	RemoveRepository() error
	RenderExists() (bool, string)
	GetCurrentDirPath() string
	GetCurrentQuartoDirPath() string
	GetCurrentZipFilePath() string
	GetCurrentRenderDirPath() string
	GetRenderFile() ([]byte, error)
	GetFileTree() (map[string]int64, error)
}
