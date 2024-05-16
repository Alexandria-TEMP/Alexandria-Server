package services_tests

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
)

var (
	versionController controllers.VersionController
	responseRecorder  *httptest.ResponseRecorder
	router            *gin.Engine
	Filesystem        filesystem.Filesystem
)

func TestMain(m *testing.M) {
	router = gin.Default()
	gin.SetMode(gin.TestMode)
	router.POST("/api/v1/version/:postID", func(c *gin.Context) {
		versionController.CreateVersion(c)
	})

	os.Exit(m.Run())
}
