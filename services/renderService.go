package services

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
	"gopkg.in/yaml.v3"
)

type RenderService struct {
	BranchRepository      database.ModelRepositoryInterface[*models.Branch]
	ProjectPostRepository database.ModelRepositoryInterface[*models.ProjectPost]
	Filesystem            filesystemInterfaces.Filesystem
}

func (renderService *RenderService) GetRenderFile(branchID uint) (string, error, error) {
	var filePath string

	// get branch
	branch, err := renderService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, nil, fmt.Errorf("failed to find branch with id %v", branchID)
	}

	// get project post
	projectPost, err := renderService.ProjectPostRepository.GetByID(branch.ProjectPostID)

	if err != nil {
		return filePath, nil, fmt.Errorf("failed to find project post with id %v", branch.ProjectPostID)
	}

	// if render is pending return 202
	if branch.RenderStatus == models.Pending {
		return filePath, fmt.Errorf("render is still pending"), nil
	}

	// if render is failed return 404
	if branch.RenderStatus == models.Failure {
		return filePath, fmt.Errorf("render has failed"), nil
	}

	// select repository of the parent post
	renderService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := renderService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		return filePath, nil, fmt.Errorf("failed to find this git branch, with name %v", branchID)
	}

	// verify render exists. if it doesn't set render status to failed
	exists, fileName := renderService.Filesystem.RenderExists()

	if !exists {
		branch.RenderStatus = models.Failure
		_, _ = renderService.BranchRepository.Update(branch)

		return filePath, nil, fmt.Errorf("render has failed")
	}

	// set filepath to absolute path to render file
	filePath = filepath.Join(renderService.Filesystem.GetCurrentRenderDirPath(), fileName)

	return filePath, nil, nil
}

func (renderService *RenderService) Render(branch *models.Branch) {
	// Unzip saved file
	if err := renderService.Filesystem.Unzip(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	// Validate project
	if valid := renderService.IsValidProject(); !valid {
		renderService.FailBranch(branch)

		return
	}

	// Install dependencies
	if err := renderService.InstallRenderDependencies(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	if err := renderService.SetProjectConfig(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	// Render quarto project
	if err := renderService.RunRender(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	// Verify that a render was produced in the form of a single file
	if exists, _ := renderService.Filesystem.RenderExists(); !exists {
		renderService.FailBranch(branch)

		return
	}

	branch.RenderStatus = models.Success
	if _, err := renderService.BranchRepository.Update(branch); err != nil {
		renderService.FailBranch(branch)

		return
	}
}

// IsValidProject validates that the files are a valid default quarto project
// They mus have a _quarto.yml or _quarto.yaml file.
// They must be of default quarto project type.
func (renderService *RenderService) IsValidProject() bool {
	// If there is no yml file to cofigure the project it is invalid
	ymlPath := filepath.Join(renderService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
	yamlPath := filepath.Join(renderService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")

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

// InstallRenderDependencies first checks if a renv.lock file is present and if so gets all dependencies.
// Next it ensures packages necessary for quarto are there.
func (renderService *RenderService) InstallRenderDependencies() error {
	// Check if renv.lock exists and if so get dependencies
	rLockPath := filepath.Join(renderService.Filesystem.GetCurrentQuartoDirPath(), "renv.lock")
	if utils.FileExists(rLockPath) {
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = renderService.Filesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	// Check if any renv exists.
	renvActivatePath := filepath.Join(renderService.Filesystem.GetCurrentQuartoDirPath(), "renv", "activate.R")
	if utils.FileExists(renvActivatePath) {
		// Install rmarkdown
		cmd := exec.Command("Rscript", "-e", "renv::install('rmarkdown')")
		cmd.Dir = renderService.Filesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}

		// Install knitr
		cmd = exec.Command("Rscript", "-e", "renv::install('knitr')")
		cmd.Dir = renderService.Filesystem.GetCurrentQuartoDirPath()
		out, err = cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(out))
		}
	}

	return nil
}

func (renderService *RenderService) SetProjectConfig() error {
	// Find config file
	yamlFilepath := filepath.Join(renderService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")
	ymlFilepath := filepath.Join(renderService.Filesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
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

	var permMode fs.FileMode = 0o666
	err = os.WriteFile(configFilepath, yamlFile, permMode)

	if err != nil {
		return fmt.Errorf("failed to write yaml config file back")
	}

	return nil
}

// RunRender renders the current project files.
// It first tries to get all dependencies and then renders to html.
func (renderService *RenderService) RunRender() error {
	// Run render command
	cmd := exec.Command("quarto", "render", renderService.Filesystem.GetCurrentQuartoDirPath(),
		"--output-dir", renderService.Filesystem.GetCurrentRenderDirPath(),
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

func (renderService *RenderService) FailBranch(branch *models.Branch) {
	branch.RenderStatus = models.Failure
	_, _ = renderService.BranchRepository.Update(branch)
}
