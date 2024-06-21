package filesystem

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/gofrs/flock"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
)

type Filesystem struct {
	CurrentDirPath       string
	CurrentQuartoDirPath string
	CurrentZipFilePath   string
	CurrentRenderDirPath string
	CurrentRepository    *git.Repository
}

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

func (filesystem *Filesystem) SetCurrentDirPath(path string) {
	filesystem.CurrentDirPath = path
}

func (filesystem *Filesystem) SetCurrentQuartoDirPath(path string) {
	filesystem.CurrentDirPath = path
}

func (filesystem *Filesystem) SetCurrentZipFilePath(path string) {
	filesystem.CurrentDirPath = path
}

func (filesystem *Filesystem) SetCurrentRenderDirPath(path string) {
	filesystem.CurrentDirPath = path
}

func InitializeFilesystem() {
	cwd, _ := os.Getwd()

	if os.Mkdir(filepath.Join(cwd, "vfs"), fs.ModePerm) != nil {
		panic("FAILED TO INITIALIZE VFS")
	}
}

func CheckoutDirectory(postID uint) interfaces.Filesystem {
	// get filepath to lock file
	cwd, _ := os.Getwd()
	rootDir := filepath.Join(cwd, "vfs")
	dirPath := filepath.Join(rootDir, strconv.FormatUint(uint64(postID), 10), "repository")

	directoryFilesystem := &Filesystem{
		CurrentDirPath:       dirPath,
		CurrentQuartoDirPath: filepath.Join(dirPath, "quarto_project"),
		CurrentZipFilePath:   filepath.Join(dirPath, "quarto_project.zip"),
		CurrentRenderDirPath: filepath.Join(dirPath, "render"),
	}

	// try to open repository if it exists.
	// we ignore the error to be flexible: if the repo already exists check it out, if not thats also ok.
	repo, _ := directoryFilesystem.CheckoutRepository()
	directoryFilesystem.CurrentRepository = repo

	return directoryFilesystem
}

func (filesystem *Filesystem) SaveZipFile(c *gin.Context, file *multipart.FileHeader) error {
	// Save zip file
	err := c.SaveUploadedFile(file, filesystem.CurrentZipFilePath)

	if err != nil {
		return fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return nil
}

// Unzip will unzip the quarto_project.zip file, if present.
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
				return fmt.Errorf("failed to make directory: %w", err)
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to make file: %w", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to open file to copy: %w", err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

// RemoveRepository entirely removes a repository
func (filesystem *Filesystem) DeleteRepository() error {
	err := os.RemoveAll(filesystem.CurrentDirPath)

	if err != nil {
		return err
	}

	return nil
}

// RenderExists checks if the render exists and is a single html file
// Returns name of the file if it exists, error if not
func (filesystem *Filesystem) RenderExists() (string, error) {
	files, err := os.ReadDir(filesystem.CurrentRenderDirPath)

	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	// Check directory contains 1 file exactly
	if len(files) != 1 {
		return "", fmt.Errorf("the directory does not contain exactly 1 file! found %d files", len(files))
	}

	// Get filename and check extension is html
	fileName := files[0].Name()

	if ext := path.Ext(fileName); ext != ".html" {
		return "", fmt.Errorf("extension '%s' is not '.html'", ext)
	}

	return fileName, nil
}

// GetFileTree returns a map of all filepaths in a quarto project and their size in bytes
func (filesystem *Filesystem) GetFileTree() (map[string]int64, error) {
	fileTree := make(map[string]int64)

	// Recursively find all files in quarto project and add path and size to map
	err := filepath.Walk(filesystem.CurrentQuartoDirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relativePath, err := filepath.Rel(filesystem.CurrentQuartoDirPath, path)

			if err != nil {
				return fmt.Errorf("the rquested file lies outside of the repository: %w", err)
			}

			// If its a directory add it with size -1
			if info.IsDir() {
				fileTree[relativePath] = -1
				return nil
			}

			fileTree[relativePath] = info.Size()

			return nil
		})

	if err != nil {
		return nil, fmt.Errorf("failed to recursively walk files: %w", err)
	}

	return fileTree, nil
}

func LockDirectory(postID uint) (*flock.Flock, error) {
	// get filepath to lock file
	cwd, _ := os.Getwd()
	lockDirPath := filepath.Join(cwd, "vfs", strconv.FormatUint(uint64(postID), 10))
	lockFilePath := filepath.Join(cwd, "vfs", strconv.FormatUint(uint64(postID), 10), "alexandria.lock")

	// check if the directory to lock exists
	if _, err := os.Stat(lockDirPath); errors.Is(err, os.ErrNotExist) {
		// create lock dir if doesn't exist
		if err := os.Mkdir(lockDirPath, fs.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create lockdir: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check if directory to lock exists: %w", err)
	}

	// check if the lockfile exists
	if _, err := os.Stat(lockFilePath); errors.Is(err, os.ErrNotExist) {
		// create lock file if doesn't exist
		if _, err := os.Create(lockFilePath); err != nil {
			return nil, fmt.Errorf("failed to create lockfile: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check if lockfile exists: %w", err)
	}

	lock := flock.New(lockFilePath)

	if err := lock.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return lock, nil
}
