package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	filesysteminterface "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type VersionService struct {
	VersionRepository database.ModelRepository[*models.Version]
	Filesystem        filesysteminterface.Filesystem
}

func (versionService *VersionService) CreateVersion(c *gin.Context, file *multipart.FileHeader, postID uint) (*models.Version, error) {
	// Create version, with pending render status
	version := models.Version{
		RenderStatus: models.Pending,
	}
	_ = versionService.VersionRepository.Create(&version)
	versionID := version.ID

	// Set paths in filesystem
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Save zip file
	if err := versionService.Filesystem.SaveRepository(c, file); err != nil {
		version.RenderStatus = models.Failure
		_, _ = versionService.VersionRepository.Update(&version)
		_ = versionService.Filesystem.RemoveRepository()

		return &version, err
	}

	// Start goroutine to render the repository.
	// This runs parallel to our response being sent, and will likely finish at a later point in time.
	// If it fails at any point we update the render status to failure and remove this repository.
	// If it succeeds we will update the renderstatus to success.
	go func() {
		// Unzip saved file
		if err := versionService.Filesystem.Unzip(); err != nil {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Validate project
		if valid := versionService.IsValidProject(); !valid {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Install dependencies
		if err := versionService.InstallRenderDependencies(); err != nil {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Render quarto project
		if err := versionService.RenderProject(); err != nil {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Verify that a render was produced in the form of a single file
		if exists, _ := versionService.Filesystem.RenderExists(); !exists {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)
			_ = versionService.Filesystem.RemoveRepository()
		}

		version.RenderStatus = models.Success
		versionService.VersionRepository.Update(&version)
	}()

	return &version, nil
}

// RenderProject renders the current project files.
// It first tries to get all dependencies and then renders to html.
func (versionService *VersionService) RenderProject() error {
	// TODO: This is super unsafe right now
	cmd := exec.Command("quarto", "render", versionService.Filesystem.GetCurrentQuartoDirPath(),
		"--output-dir", versionService.Filesystem.GetCurrentRenderDirPath(),
		"--to", "html",
		"--no-cache",
		"-M", "embed-resources:true",
		"-M", "toc-location:body",
		"-M", "margin-left:0",
		"-M", "margin-right:0",
		"--log-level", "error",
	)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

// InstallRenderDependencies first checks if a renv.lock file is present and if so gets all dependencies.
// Next it ensures packages necessary for quarto are there.
func (versionService *VersionService) InstallRenderDependencies() error {
	// Check if renv.lock exists and if so get dependencies
	rLockPath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "renv.lock")
	if filesystem.FileExists(rLockPath) {
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = versionService.Filesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	// Check if any renv exists.
	renvActivatePath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "renv", "activate.R")
	if filesystem.FileExists(renvActivatePath) {
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

	if !(filesystem.FileExists(ymlPath)) {
		ymlPath = yamlPath

		if !(filesystem.FileExists(yamlPath)) {
			return false
		}
	}

	// If they type is not default it is invalid
	if filesystem.FileContains(ymlPath, "type:") && !filesystem.FileContains(ymlPath, "type: default") {
		return false
	}

	return true
}

func (versionService *VersionService) GetRenderFile(versionID, postID uint) (string, error, error) {
	version, err := versionService.VersionRepository.GetByID(versionID)
	var filePath string

	if err != nil {
		return filePath, nil, fmt.Errorf("no such version exists")
	}

	// If pending return error 202
	if version.RenderStatus == models.Pending {
		return filePath, fmt.Errorf("version still rendering"), nil
	}

	// If failure return error 404
	if version.RenderStatus == models.Failure {
		return filePath, nil, fmt.Errorf("version failed to render")
	}

	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Check that render exists, if not update render status to failed and return 404
	if exists, _ := versionService.Filesystem.RenderExists(); !exists {
		version.RenderStatus = models.Failure
		_, _ = versionService.VersionRepository.Update(version)

		return filePath, nil, fmt.Errorf("version failed to render")
	}

	return versionService.Filesystem.GetCurrentRenderDirPath(), nil, nil
}

func (versionService *VersionService) GetRepositoryFile(versionID, postID uint) (string, error) {
	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Check that render exists, if not update render status to failed and return 404
	if exists := filesystem.FileExists(versionService.Filesystem.GetCurrentZipFilePath()); !exists {
		return "", fmt.Errorf("no such file exists")
	}

	return versionService.Filesystem.GetCurrentZipFilePath(), nil
}
