package servicestests

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gorm.io/gorm"
)

var (
	versionService services.VersionService
	c              *gin.Context
	mockFilesystem *mocks.MockFilesystem
	exampleVersion models.Version
	cwd            string
)

func TestMain(m *testing.M) {
	exampleVersion = models.Version{
		Model:        gorm.Model{ID: 0},
		Discussions:  nil,
		RenderStatus: models.Pending,
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
