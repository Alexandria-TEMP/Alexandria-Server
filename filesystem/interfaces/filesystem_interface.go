package filesystem_interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	SetCurrentVersion(versionID, postID uint)
	SaveRepository(c *gin.Context, file *multipart.FileHeader) error
	Unzip() error
	RemoveProjectDirectory() error
	RemoveRepository() error
	CountRenderFiles() int
	GetCurrentDirPath() string
	GetCurrentQuartoDirPath() string
	GetCurrentZipFilePath() string
	GetCurrentRenderDirPath() string
	GetRenderFile() (forms.OutgoingFileForm, string, error)
	GetRepositoryFile() (forms.OutgoingFileForm, string, error)
}
