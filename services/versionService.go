package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesysteminterface "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
	"gopkg.in/yaml.v3"
)

type VersionService struct {
	VersionRepository database.RepositoryInterface[*models.Version]
	Filesystem        filesysteminterface.Filesystem
}

func (versionService *VersionService) CreateVersion(c *gin.Context, file *multipart.FileHeader, postID uint) (*models.Version, error) {
	// Create version, with pending render status
	version := models.Version{RenderStatus: models.Pending}
	_ = versionService.VersionRepository.Create(&version)

	// Set paths in filesystem
	versionService.Filesystem.SetCurrentVersion(version.ID, postID)

	// Save zip file
	if err := versionService.Filesystem.SaveRepository(c, file); err != nil {
		version.RenderStatus = models.Failure
		_, _ = versionService.VersionRepository.Update(&version)
		_ = versionService.Filesystem.RemoveRepository()

		return &version, err
	}

	// This goroutine runs parallel to our response being sent, and will likely finish at a later point in time.
	// If it fails before rendering we will remove the repository entirely.
	// If it succeeds we will update the renderstatus to success, otherwise failure.
	go func() {
		// Unzip saved file
		if err := versionService.Filesystem.Unzip(); err != nil {
			versionService.FailAndRemoveVersion(&version)

			return
		}

		// Validate project
		if valid := versionService.IsValidProject(); !valid {
			versionService.FailAndRemoveVersion(&version)

			return
		}

		// Install dependencies
		if err := versionService.InstallRenderDependencies(); err != nil {
			versionService.FailAndRemoveVersion(&version)

			return
		}

		if err := versionService.SetProjectConfig(); err != nil {
			versionService.FailAndRemoveVersion(&version)

			return
		}

		// Render quarto project
		if err := versionService.RenderProject(); err != nil {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)

			return
		}

		// Verify that a render was produced in the form of a single file
		if exists, _ := versionService.Filesystem.RenderExists(); !exists {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)

			return
		}

		version.RenderStatus = models.Success
		if _, err := versionService.VersionRepository.Update(&version); err != nil {
			version.RenderStatus = models.Failure
			_, _ = versionService.VersionRepository.Update(&version)

			return
		}
	}()

	return &version, nil
}

func (versionService *VersionService) SetProjectConfig() error {
	// Find config file
	yamlFilepath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")
	ymlFilepath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
	configFilepath := yamlFilepath

	if !utils.FileExists(yamlFilepath) {
		configFilepath = ymlFilepath
	}

	// Unmarshal yaml file
	yamlObj := make(map[string]interface{})
	yamlFile, err := os.ReadFile(configFilepath)

	if err != nil {
		return fmt.Errorf("failed to open yaml config file")
	}

	err = yaml.Unmarshal(yamlFile, yamlObj)

	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml config file")
	}

	yamlObj["format"] = map[string]interface{}{"html": map[string]interface{}{"page-layout": "custom"}}
	yamlFile, err = yaml.Marshal(yamlObj)

	if err != nil {
		return fmt.Errorf("failed to marshal yaml config file")
	}

	err = os.WriteFile(configFilepath, yamlFile, 0666)

	if err != nil {
		return fmt.Errorf("failed to write yaml config file back")
	}

	return nil
}

// RenderProject renders the current project files.
// It first tries to get all dependencies and then renders to html.
func (versionService *VersionService) RenderProject() error {
	// Run render command
	cmd := exec.Command("quarto", "render", versionService.Filesystem.GetCurrentQuartoDirPath(),
		"--output-dir", versionService.Filesystem.GetCurrentRenderDirPath(),
		"--to", "html",
		"--no-cache",
		"-M", "embed-resources:true",
		"-M", "title:",
		"-M", "date:",
		"-M", "date-modified:",
		"-M", "author:",
		"-M", "doi:",
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
	if utils.FileExists(rLockPath) {
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = versionService.Filesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	// Check if any renv exists.
	renvActivatePath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "renv", "activate.R")
	if utils.FileExists(renvActivatePath) {
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
// They mus have a _quarto.yml or _quarto.yaml file.
// They must be of default quarto project type.
func (versionService *VersionService) IsValidProject() bool {
	// If there is no yml file to cofigure the project it is invalid
	ymlPath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
	yamlPath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")

	if !(utils.FileExists(ymlPath)) {
		ymlPath = yamlPath

		if !(utils.FileExists(yamlPath)) {
			return false
		}
	}

	// If they type is not default it is invalid
	if utils.FileContains(ymlPath, "type:") && !utils.FileContains(ymlPath, "type: default") {
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
	if exists := utils.FileExists(versionService.Filesystem.GetCurrentZipFilePath()); !exists {
		return "", fmt.Errorf("no such file exists")
	}

	absFilepath, _ := filepath.Abs(versionService.Filesystem.GetCurrentZipFilePath())

	return absFilepath, nil
}

func (versionService *VersionService) GetTreeFromRepository(versionID, postID uint) (map[string]int64, error, error) {
	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)

	// Check that render exists, if not update render status to failed and return 404
	if exists := utils.FileExists(versionService.Filesystem.GetCurrentQuartoDirPath()); !exists {
		return nil, fmt.Errorf("no such directory exists"), nil
	}

	fileTree, err := versionService.Filesystem.GetFileTree()

	return fileTree, nil, err
}

func (versionService *VersionService) GetFileFromRepository(versionID, postID uint, relFilepath string) (string, error) {
	if strings.Contains(relFilepath, "..") {
		return "", fmt.Errorf("file is outside of repository")
	}

	// Set current version
	versionService.Filesystem.SetCurrentVersion(versionID, postID)
	absFilepath := filepath.Join(versionService.Filesystem.GetCurrentQuartoDirPath(), relFilepath)

	// Check that file exists, if not return 404
	if exists := utils.FileExists(absFilepath); !exists {
		return "", fmt.Errorf("no such file exists")
	}

	return absFilepath, nil
}

func (versionService *VersionService) FailAndRemoveVersion(version *models.Version) {
	version.RenderStatus = models.Failure
	_, _ = versionService.VersionRepository.Update(version)
	_ = versionService.Filesystem.RemoveRepository()
}
