package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./versionService_interface.go -destination=../../mocks/versionService_mock.go

type VersionService interface {
	CreateVersion(c *gin.Context, file *multipart.FileHeader, postID uint) (*models.Version, error)
	RenderProject() error
	GetRender(versionID, postID uint) (forms.OutgoingFileForm, string, error)
	GetRepository(versionID, postID uint) (forms.OutgoingFileForm, string, error)
}
