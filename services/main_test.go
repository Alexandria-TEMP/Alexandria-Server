package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	versionService VersionService
	c              *gin.Context
	mockFilesystem *mocks.MockFilesystem
	exampleVersion models.Version
	cwd            string
)

func TestMain(m *testing.M) {
	exampleVersion = models.Version{
		Model:        gorm.Model{ID: 0},
		Discussions:  nil,
		RenderStatus: models.RenderPending,
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
