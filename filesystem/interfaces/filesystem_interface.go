package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

//go:generate mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go

type Filesystem interface {
	CheckoutDirectory(postID uint)
	GetFileTree() (map[string]int64, error)

	CreateRepository() error
	CheckoutRepository() (*git.Repository, error)
	DeleteRepository() error

	CreateBranch(branchName string) error
	CheckoutBranch(branchName string) error

	CreateCommit() error
	GetLastCommit() (*plumbing.Reference, error)

	SaveZipFile(c *gin.Context, file *multipart.FileHeader) error
	Unzip() error

	RenderExists() (bool, string)

	GetCurrentDirPath() string
	GetCurrentQuartoDirPath() string
	GetCurrentZipFilePath() string
	GetCurrentRenderDirPath() string
}
