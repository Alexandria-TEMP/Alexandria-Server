package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	// CheckoutDirectory set filepaths according to postID.
	// If a git repo exists there it will be opened.
	// CurrentDirPath = <cwd>/vfs/<postID>
	// CurrentQuartoDirPath = <cwd>/vfs/<postID>/quarto_project
	// CurrentZipFilePath = <cwd>/vfs/<postID>/quarto_project.zip
	// CurrentRenderDirPath = <cwd>/vfs/<postID>/render/<some_html_file>
	CheckoutDirectory(postID uint)

	// CreateRepository create git repository at CurrentDirPath
	CreateRepository() error

	// CheckoutRepository checkout git repository if there is one at CurrentDirPath
	CheckoutRepository() (*git.Repository, error)

	// DeleteRepository delete entire repository and directory at CurrentDirPath
	DeleteRepository() error

	// CreateBranch create a new branch off of master's last commit
	CreateBranch(branchName string) error

	// DeleteBranch delete a branch
	DeleteBranch(branchName string) error

	// CheckoutBranch checkout branch
	CheckoutBranch(branchName string) error

	// CreateCommit commit current changes
	CreateCommit() error

	// GetLastCommit get last commit reference for specific branch
	GetLastCommit(branchName string) (*plumbing.Reference, error)

	// GetFileTree get all files at GetCurrentQuartoDirPath
	GetFileTree() (map[string]int64, error)

	// RenderExists checks whether render exists as expected at GetCurrentRenderDirPath and returns filename
	RenderExists() (bool, string)

	// SaveZipFile saves a zip file from the gin context to GetCurrentZipFilePath
	SaveZipFile(c *gin.Context, file *multipart.FileHeader) error

	// Unzip unzips a project to GetCurrentQuartoDirPath
	Unzip() error

	// CleanDir deletes all files from the currnet branch
	CleanDir() error

	// Reset resets a branch to the last commit
	Reset() error

	GetCurrentDirPath() string
	GetCurrentQuartoDirPath() string
	GetCurrentZipFilePath() string
	GetCurrentRenderDirPath() string
}
