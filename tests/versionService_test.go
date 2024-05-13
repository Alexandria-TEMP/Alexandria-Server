package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

func TestSaveRepository(t *testing.T) {
	versionController := controllers.VersionController{
		VersionService: services.VersionService{
			Filesystem: filesystem.InitFilesystem(),
		},
	}

	body, dataType := utils.CreateMultipartFile("file.zip")

	responseRecorder = httptest.NewRecorder()

	router = gin.Default()
	gin.SetMode(gin.TestMode)
	router.POST("/api/v1/version/:postID", func(c *gin.Context) {
		versionController.CreateVersion(c)
	})

	req, _ := http.NewRequest("POST", "/api/v1/version/1", body)
	req.Header.Add("Content-Type", dataType)
	router.ServeHTTP(responseRecorder, req)

}
