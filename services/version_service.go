package services

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type VersionService struct {
	Filesystem filesystem.Filesystem
}

func (versionService VersionService) GetFilesystem() filesystem.Filesystem {
	return versionService.Filesystem
}

func (versionService VersionService) SaveRepository(c *gin.Context, file *multipart.FileHeader, versionID, postID uint) error {
	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Save zip file
	err := versionService.Filesystem.SaveRepository(c, file)

	if err != nil {
		return err
	}

	// Unzip saved file
	err = versionService.Filesystem.Unzip()

	if err != nil {
		return err
	}

	// Render quarto project
	err = versionService.Filesystem.RenderProject()

	return nil
}

func (versionService VersionService) CreateVersion() *models.Version {
	return new(models.Version)
}
