package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gofrs/flock"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	// CheckoutDirectory set filepaths according to postID.
	// If a git repo exists there it will be opened.
	// CurrentDirPath = <cwd>/vfs/<postID>
	// CurrentQuartoDirPath = <cwd>/vfs/<postID>/quarto_project
	// CurrentZipFilePath = <cwd>/vfs/<postID>/quarto_project.zip
	// CurrentRenderDirPath = <cwd>/vfs/<postID>/render/<some_html_file>
	// CheckoutDirectory(postID uint) Filesystem

	// LockDirectory locks the directory associated with the post.
	// It returns a lock, which must be unlocked after changes are done.
	// LockDirectory(postID uint) (*flock.Flock, error)

	// CreateRepository create git repository at CurrentDirPath
	CreateRepository() error

	// CheckoutRepository checkout git repository if there is one at CurrentDirPath
	CheckoutRepository() (*git.Repository, error)

	// CreateBranch create a new branch off of master's last commit
	CreateBranch(branchName string) error

	// DeleteBranch delete a branch
	DeleteBranch(branchName string) error

	// CheckoutBranch checkout branch
	CheckoutBranch(branchName string) error

	// Merge actually resets master to the last commit on the branch we are merging
	Merge(toMerge, mergeInto string) error

	// CreateCommit commit current changes
	CreateCommit() error

	// GetLastCommit get last commit reference for specific branch
	GetLastCommit(branchName string) (*plumbing.Reference, error)

	// GetFileTree get all files at GetCurrentQuartoDirPath
	GetFileTree() (map[string]int64, error)

	// RenderExists checks whether render exists as expected at GetCurrentRenderDirPath and returns filename
	RenderExists() (string, error)

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
	SetCurrentDirPath(string)
	SetCurrentQuartoDirPath(string)
	SetCurrentZipFilePath(string)
	SetCurrentRenderDirPath(string)
}

type FilesystemManagerInterface interface {
	LockDirectory(postID uint) (*flock.Flock, error)
	CheckoutDirectory(postID uint) Filesystem
}
