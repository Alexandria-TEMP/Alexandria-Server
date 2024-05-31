package services

import (
	"errors"
	"mime/multipart"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
	filesysteminterface "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type VersionService struct {
	Filesystem        filesysteminterface.Filesystem
	VersionRepository database.RepositoryInterface[*models.Version]
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
	version := models.Version{
		RenderStatus: models.Pending,
	}
	versionID := version.ID

	// Set paths in filesystem
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Save zip file
	if err := versionService.Filesystem.SaveRepository(c, file); err != nil {
		_ = versionService.Filesystem.RemoveRepository()
		return &version, err
	}

	// Start goroutine to render after responding to client.
	// If it fails at any point we update the renderstatus to failure and remove the directory to this repository.
	// If it succeeds we will update the renderstatus to success and remove the quarto project directory.
	go func() {
		// Unzip saved file
		if err := versionService.Filesystem.Unzip(); err != nil {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Validate project
		if valid := versionService.IsValidProject(); !valid {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		// Install dependencies
		if err := versionService.InstallRenderDependencies(); err != nil {
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

		// Verify that a render was produced in the form of a single file
		if numFiles := versionService.Filesystem.CountRenderFiles(); numFiles != 1 {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()
		}

		// Remove unzipped project file
		if err := versionService.Filesystem.RemoveProjectDirectory(); err != nil {
			version.RenderStatus = models.Failure
			_ = versionService.Filesystem.RemoveRepository()

			return
		}

		version.RenderStatus = models.Success
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

func (versionService *VersionService) GetRender(versionID, postID uint) (forms.OutgoingFileForm, string, error) {
	// TODO: Check version render status
	// If pending return error eccordingly
	// If failure return error accordingly
	// If success proceed to steps below
	//
	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Get render file
	return versionService.Filesystem.GetRenderFile()
}

func (versionService *VersionService) GetRepository(versionID, postID uint) (forms.OutgoingFileForm, string, error) {
	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Get repository file
	return versionService.Filesystem.GetRepositoryFile()
}
