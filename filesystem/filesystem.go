package filesystem

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
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

func (filesystem *Filesystem) GetCurrentDirPath() string {
	return filesystem.CurrentDirPath
}

func (filesystem *Filesystem) GetCurrentQuartoDirPath() string {
	return filesystem.CurrentQuartoDirPath
}

func (filesystem *Filesystem) GetCurrentZipFilePath() string {
	return filesystem.CurrentZipFilePath
}

func (filesystem *Filesystem) GetCurrentRenderDirPath() string {
	return filesystem.CurrentRenderDirPath
}

// InitFilesystem initializes a new filesystem by setting the root to the current working directory and assigning default values.
func InitFilesystem() *Filesystem {
	filesystem := &Filesystem{
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

// SetCurrentVersion will set the paths the filesystem uses in accordance with the IDs passed.
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

// RemoveRepository entirely removes a version repository if it is not valid
func (filesystem *Filesystem) RemoveRepository() error {
	err := os.RemoveAll(filesystem.CurrentDirPath)

	if err != nil {
		return err
	}

	return nil
}

// CountRenderFiles counts how many files are at the render directory of this version
func (filesystem *Filesystem) CountRenderFiles() int {
	files, err := os.ReadDir(filesystem.CurrentRenderDirPath)

	if err != nil {
		return 0
	}

	return len(files)
}

// Returns the rendered project post wrapped in a OutgoingFileForm with the content-type
func (filesystem *Filesystem) GetRenderFile() (forms.OutgoingFileForm, string, error) {
	var outgoingFileForm forms.OutgoingFileForm

	// Check if directory exists
	files, err := os.ReadDir(filesystem.CurrentRenderDirPath)

	if err != nil {
		return outgoingFileForm, "", errors.New("invalid directory")
	}

	// Check directory contains 1 file exactly
	if len(files) != 1 {
		return outgoingFileForm, "", fmt.Errorf("expected 1 file, but found %v", len(files))
	}

	// Get filename and check extension is html
	fileName := files[0].Name()

	if ext := path.Ext(fileName); ext != ".html" {
		return outgoingFileForm, "", fmt.Errorf("expected .html file, but found %v", ext)
	}

	// Create multipart file
	filePath := filepath.Join(filesystem.CurrentRenderDirPath, fileName)
	rawFile, contentType, err := CreateMultipartFile(filePath)

	if err != nil {
		return outgoingFileForm, "", err
	}

	outgoingFileForm = forms.OutgoingFileForm{
		File: &rawFile,
	}

	return outgoingFileForm, contentType, nil
}

// GetRepositoryFile return as zipped quarto project after validating that it exists, together the content type
func (filesystem *Filesystem) GetRepositoryFile() (forms.OutgoingFileForm, string, error) {
	var outgoingFileForm forms.OutgoingFileForm

	// Check if file exists
	if !FileExists(filesystem.CurrentZipFilePath) {
		return outgoingFileForm, "", errors.New("this project doesn't exist")
	}

	// Create multipart file
	rawFile, contentType, err := CreateMultipartFile(filesystem.CurrentZipFilePath)

	if err != nil {
		return outgoingFileForm, "", err
	}

	outgoingFileForm = forms.OutgoingFileForm{
		File: &rawFile,
	}

	return outgoingFileForm, contentType, nil
}
