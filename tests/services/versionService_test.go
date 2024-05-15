package services_tests

import (
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/controllers"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services"
)

func beforeEach(t *testing.T) {
	t.Helper()

	versionController = controllers.VersionController{
		VersionService: services.VersionService{
			Filesystem: filesystem.InitFilesystem(),
		},
	}

	responseRecorder = httptest.NewRecorder()
}

func cleanup(t *testing.T) {
	t.Helper()

	if cwd, err := os.Getwd(); err != nil {
		os.RemoveAll(filepath.Join(cwd, "vfs"))
	}
}

// func TestSaveRepository200(t *testing.T) {
// 	beforeEach(t)

// 	body, dataType := filesystem.CreateMultipartFile("file.zip")

// 	req, _ := http.NewRequest("POST", "/api/v1/version/1", body)
// 	req.Header.Add("Content-Type", dataType)
// 	router.ServeHTTP(responseRecorder, req)

// 	defer responseRecorder.Result().Body.Close()

// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)

// 	// cleanup(t)
// }

func TestSetCurrentVersion(t *testing.T) {
	beforeEach(t)

	Filesystem := versionController.VersionService.GetFilesystem()

	Filesystem.SetCurrentVersion(5, 10)

	cwd, _ := os.Getwd()
	assert.Equal(t, filepath.Join(cwd, "vfs", "10", "5"), Filesystem.CurrentDirPath)
	assert.Equal(t, filepath.Join(cwd, "vfs", "10", "5", "quarto_project"), Filesystem.CurrentQuartoDirPath)
	assert.Equal(t, filepath.Join(cwd, "vfs", "10", "5", "render"), Filesystem.CurrentRenderDirPath)
	assert.Equal(t, filepath.Join(cwd, "vfs", "10", "5", "quarto_project.zip"), Filesystem.CurrentZipFilePath)
}

func TestUnzipSuccess(t *testing.T) {
	beforeEach(t)

	Filesystem := versionController.VersionService.GetFilesystem()
	Filesystem.SetCurrentVersion(0, 0)
	err := Filesystem.Unzip()

	if err != nil {
		println(err)
	}

	cleanup(t)
}

func TestRenderSuccess(t *testing.T) {
	beforeEach(t)

	Filesystem := versionController.VersionService.GetFilesystem()
	Filesystem.SetCurrentVersion(0, 0)
	err := Filesystem.RenderProject()

	if err != nil {
		fmt.Printf("%v", err)
	}
}
