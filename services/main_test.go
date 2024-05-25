package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

var (
	versionService    VersionService
	c                 *gin.Context
	mockFilesystem    *mocks.MockFilesystem
	pendingVersion    models.Version
	failureVersion    models.Version
	successVersion    models.Version
	cwd               string
	versionRepository database.ModelRepository[*models.Version]
	db                *gorm.DB
)

func TestMain(m *testing.M) {
	pendingVersion = models.Version{
		Model:        gorm.Model{ID: 0},
		Discussions:  nil,
		RenderStatus: models.Pending,
	}

	failureVersion = models.Version{
		Model:        gorm.Model{ID: 1},
		Discussions:  nil,
		RenderStatus: models.Failure,
	}

	successVersion = models.Version{
		Model:        gorm.Model{ID: 2},
		Discussions:  nil,
		RenderStatus: models.Success,
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
