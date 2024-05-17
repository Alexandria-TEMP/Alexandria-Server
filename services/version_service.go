package services

import (
	"errors"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	filesystem_interfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type VersionService struct {
	Filesystem filesystem_interfaces.Filesystem
}

// GetFilesystem is a helper function to test the version service
func (versionService *VersionService) GetFilesystem() *filesystem_interfaces.Filesystem {
	return &versionService.Filesystem
}

// CreateVersion orchestrates the version creation.
// 1. creates a new version with the pending render status.
// 2. saves the file to its directory
// 3. return a status 200 to client
// 4. unzip file
// 5. render project and update render status
// 6. delete unzipped project files
// TODO: persist data
func (versionService *VersionService) CreateVersion(c *gin.Context, file *multipart.FileHeader, postID uint) (*models.Version, error) {
	version := &models.Version{
		RenderStatus: models.Pending,
	}
	versionID := version.ID

	// Set paths in filesystem
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Save zip file
	if err := versionService.Filesystem.SaveRepository(c, file); err != nil {
		_ = versionService.Filesystem.RemoveRepository()
		return version, err
	}

	// Start goroutine to render after responding to client
	go func() {
		// Unzip saved file
		if err := versionService.Filesystem.Unzip(); err != nil {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		if valid := versionService.IsValidProject(); !valid {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Render quarto project
		if err := versionService.RenderProject(); err != nil {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Remove unzipped project file
		if err := versionService.Filesystem.RemoveProjectDirectory(); err != nil {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		version.RenderStatus = models.Success
	}()

	return version, nil
}

// RenderProject renders the current project files.
// It first tries to get all dependencies and then renders to html.
func (versionService *VersionService) RenderProject() error {
	err := versionService.installRenderDependencies()

	if err != nil {
		return err
	}

	// TODO: This is super unsafe right now
	cmd := exec.Command("quarto", "render", versionService.Filesystem.GetCurrentQuartoDirPath(),
		"--output-dir", versionService.Filesystem.GetCurrentRenderDirPath(),
		"--to", "html",
		"--no-cache",
		"-M", "embed-resources:true",
		"-M", "toc-location:body")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

// InstallRenderDependencies first checks if a renv.lock file is present and if so gets all dependencies.
// Next it ensures packages necessary for quarto are there.
func (versionService *VersionService) installRenderDependencies() error {
	// Check if renv.lock exists and if so get dependencies
	rLockPath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "renv.lock")
	if FileExists(rLockPath) {
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = versionService.Filesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	// Check if any renv exists.
	renvActivatePath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "renv", "activate.R")
	if FileExists(renvActivatePath) {
		// Install rmarkdown
		cmd := exec.Command("Rscript", "-e", "renv::install('rmarkdown')")
		cmd.Dir = versionService.Filesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}

		// Install knitr
		cmd = exec.Command("Rscript", "-e", "renv::install('knitr')")
		cmd.Dir = versionService.Filesystem.GetCurrentQuartoDirPath()
		out, err = cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	return nil
}

// IsValidProject validates that the files are a valid default quarto project
func (versionService *VersionService) IsValidProject() bool {
	// If there is no yml file to cofigure the project it is invalid
	ymlPath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
	yamlPath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")

	if !(FileExists(ymlPath)) {
		ymlPath = yamlPath

		if !(FileExists(yamlPath)) {
			return false
		}
	}

	// If they type is not default it is invalid
	if FileContains(ymlPath, "type:") && !FileContains(ymlPath, "type: default") {
		return false
	}

	return true
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func FileContains(filePath, match string) bool {
	if !FileExists(filePath) {
		return false
	}

	text, _ := os.ReadFile(filePath)

	return strings.Contains(string(text), match)
}
