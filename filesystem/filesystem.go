package filesystem

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
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

// RemoveRepository entirely removes a version repository if it is not valid
func (filesystem *Filesystem) RemoveRepository() error {
	err := os.RemoveAll(filesystem.CurrentDirPath)

	if err != nil {
		return err
	}

	return nil
}

// RenderExists checks if the render exists and is a single html file
// Returns a bool and the name of the file if it does exist
func (filesystem *Filesystem) RenderExists() (exists bool, name string) {
	files, err := os.ReadDir(filesystem.CurrentRenderDirPath)

	if err != nil {
		return false, ""
	}

	// Check directory contains 1 file exactly
	if len(files) != 1 {
		return false, ""
	}

	// Get filename and check extension is html
	fileName := files[0].Name()

	if ext := path.Ext(fileName); ext != ".html" {
		return false, ""
	}

	return true, fileName
}

// Returns the rendered project as binary large object, ie a byte slice
func (filesystem *Filesystem) GetRenderFile() ([]byte, error) {
	// Check if directory exists
	exists, fileName := filesystem.RenderExists()

	if !exists {
		return nil, fmt.Errorf("render doesn't exist or is invalid")
	}

	// Create blob
	filePath := filepath.Join(filesystem.CurrentRenderDirPath, fileName)
	file, err := CreateByteSliceFile(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetFileTree returns a map of all filepaths in a quarto project and their size in bytes
func (filesystem Filesystem) GetFileTree() (map[string]int64, error) {
	var fileTree map[string]int64

	// Recursively find all files in quarto project and add path and size to map
	err := filepath.Walk(filesystem.CurrentQuartoDirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relativePath, err := filepath.Rel(filesystem.CurrentQuartoDirPath, path)

			if err != nil {
				return err
			}

			fileTree[relativePath] = info.Size()

			return nil
		})

	if err != nil {
		return nil, nil
	}

	return fileTree, nil
}
