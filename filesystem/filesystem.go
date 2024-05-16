package filesystem

import (
	"archive/zip"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Filesystem struct {
	rootPath             string
	zipName              string
	quartoDirectoryName  string
	CurrentDirPath       string
	CurrentQuartoDirPath string
	CurrentZipFilePath   string
	CurrentRenderDirPath string
}

var (
	cwd, _                     = os.Getwd()
	defaultRootPath            = filepath.Clean(filepath.Join(cwd, "vfs"))
	defaultZipName             = "quarto_project.zip"
	defaultQuartoDirectoryName = "quarto_project"
)

func InitFilesystem() Filesystem {
	filesystem := Filesystem{
		rootPath:            defaultRootPath,
		zipName:             defaultZipName,
		quartoDirectoryName: defaultQuartoDirectoryName,
	}

	err := os.MkdirAll(filesystem.rootPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return filesystem
}

func (filesystem *Filesystem) SetCurrentVersion(versionID, postID uint) {
	filesystem.CurrentDirPath = filepath.Join(filesystem.rootPath, strconv.FormatUint(uint64(postID), 10), strconv.FormatUint(uint64(versionID), 10))
	filesystem.CurrentQuartoDirPath = filepath.Join(filesystem.CurrentDirPath, filesystem.quartoDirectoryName)
	filesystem.CurrentZipFilePath = filepath.Join(filesystem.CurrentDirPath, filesystem.zipName)
	filesystem.CurrentRenderDirPath = filepath.Join(filesystem.CurrentDirPath, "render")
}

// SaveRepository saves a zip file to a ./vfs/{postID}/{versionID} in the filesystem and return the path to the directory.
func (filesystem *Filesystem) SaveRepository(c *gin.Context, file *multipart.FileHeader) error {
	// Save zip file
	err := c.SaveUploadedFile(file, filesystem.CurrentZipFilePath)

	if err != nil {
		return err
	}

	return nil
}

// Unzip will unzip the quarto_project.zip file, if present, of any post version.
// Errors if there is no such file or it can't unzip it.
func (filesystem *Filesystem) Unzip() error {
	archive, err := zip.OpenReader(filesystem.CurrentZipFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(filesystem.CurrentQuartoDirPath, f.Name)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)

			if err != nil {
				return err
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

// RemoveProjectDirectory only removes the unzipped files, not the zip file or the render
func (filesystem *Filesystem) RemoveProjectDirectory() error {
	err := os.RemoveAll(filesystem.CurrentQuartoDirPath)

	if err != nil {
		return err
	}

	return nil
}
