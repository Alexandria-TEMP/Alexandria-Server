package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

type VersionService struct {
	Filesystem filesystem.Filesystem
}

func (versionService VersionService) SaveRepository(c *gin.Context, file *multipart.FileHeader, versionID uint, postID uint) error {
	dirPath := versionService.Filesystem.GetRepositoryPath(versionID, postID)
	zipName := fmt.Sprintf("%s.zip", strconv.FormatUint(uint64(versionID), 10))
	zipFilePath := filepath.Join(dirPath, zipName)

	err := c.SaveUploadedFile(file, zipFilePath)

	if err != nil {
		return err
	}

	err = utils.Unzip(zipFilePath, dirPath)

	if err != nil {
		return err
	}

	return nil
}

func (versionService VersionService) CreateVersion() *models.Version {
	return new(models.Version)
}
