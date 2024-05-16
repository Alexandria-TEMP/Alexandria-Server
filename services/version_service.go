package services

import (
	"errors"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type VersionService struct {
	Filesystem filesystem.Filesystem
}

// GetFilesystem is a helper function to test the version service
func (versionService *VersionService) GetFilesystem() *filesystem.Filesystem {
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
		return version, err
	}

	// Start goroutine to render after responding to client
	go func() {
		// Unzip saved file
		if err := versionService.Filesystem.Unzip(); err != nil {
			version.RenderStatus = models.Failure
			return
		}

		// Render quarto project
		if err := versionService.RenderProject(); err != nil {
			version.RenderStatus = models.Failure
			return
		}

		// Remove unzipped project file
		if err := versionService.Filesystem.RemoveProjectDirectory(); err != nil {
			version.RenderStatus = models.Failure
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
	cmd := exec.Command("quarto", "render", versionService.Filesystem.CurrentQuartoDirPath,
		"--output-dir", versionService.Filesystem.CurrentRenderDirPath,
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
	rLockPath := filepath.Join(versionService.Filesystem.CurrentQuartoDirPath, "renv.lock")
	if _, err := os.Stat(rLockPath); err == nil {
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = versionService.Filesystem.CurrentQuartoDirPath
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	// Install rmarkdown
	cmd := exec.Command("Rscript", "-e", "renv::install('rmarkdown')")
	cmd.Dir = versionService.Filesystem.CurrentQuartoDirPath
	out, err := cmd.CombinedOutput()

	if err != nil {
		return errors.New(string(out))
	}

	// Install knitr
	cmd = exec.Command("Rscript", "-e", "renv::install('knitr')")
	cmd.Dir = versionService.Filesystem.CurrentQuartoDirPath
	out, err = cmd.CombinedOutput()

	if err != nil {
		return errors.New(string(out))
	}

	return nil
}
