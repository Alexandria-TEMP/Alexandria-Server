package interfaces

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./postService_interface.go -destination=../../mocks/postService_mock.go

type PostService interface {
	GetPost(postID uint) (*models.Post, error)
	CreatePost(form *forms.PostCreationForm) (*models.Post, error)
	UpdatePost(updatedPost *models.Post) error

	// UploadPost saves a zipped quarto project to master and initiates the render pipeline.
	// It the renders the project in a goroutine.
	UploadPost(c *gin.Context, file *multipart.FileHeader, postID uint) error
	GetMainProject(postID uint) (string, error)
	GetMainFiletree(branchID uint) (map[string]int64, error, error)
	GetMainFileFromProject(postID uint, relFilepath string) (string, error)
}
