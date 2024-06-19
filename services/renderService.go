package services

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofrs/flock"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	filesystemInterfaces "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
	"gopkg.in/yaml.v3"
)

type RenderService struct {
	BranchRepository      database.ModelRepositoryInterface[*models.Branch]
	PostRepository        database.ModelRepositoryInterface[*models.Post]
	ProjectPostRepository database.ModelRepositoryInterface[*models.ProjectPost]
	Filesystem            filesystemInterfaces.Filesystem
	BranchService         interfaces.BranchService
}

func (renderService *RenderService) GetRenderFile(branchID uint) (string, *flock.Flock, error, error) {
	var filePath string

	// get branch
	branch, err := renderService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, nil, nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := renderService.BranchService.GetBranchProjectPost(branch)

	if err != nil {
		return filePath, nil, nil, fmt.Errorf("failed to find project post with id %v: %w", branch.ProjectPostID, err)
	}

	// if render is pending return 202
	if branch.RenderStatus == models.Pending {
		return filePath, nil, fmt.Errorf("render is still pending"), nil
	}

	// if render is failed return 404
	if branch.RenderStatus == models.Failure {
		return filePath, nil, nil, fmt.Errorf("render has failed")
	}

	// lock directory
	// unlock upon error or after controller has read file
	lock, err := renderService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return filePath, nil, nil, fmt.Errorf("failed to aquire lock for directory %v: %w", projectPost.PostID, err)
	}

	// select repository of the parent post
	renderService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := renderService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		lock.Unlock()
		return filePath, nil, nil, fmt.Errorf("failed to find this git branch, with name %v: %w", branchID, err)
	}

	// verify render exists. if it doesn't set render status to failed
	fileName, err := renderService.Filesystem.RenderExists()
	if err != nil {
		branch.RenderStatus = models.Failure
		_, _ = renderService.BranchRepository.Update(branch)
		lock.Unlock()

		return filePath, nil, nil, fmt.Errorf("render has failed: %w", err)
	}

	// set filepath to absolute path to render file
	filePath = filepath.Join(renderService.Filesystem.GetCurrentRenderDirPath(), fileName)

	return filePath, lock, nil, nil
}

func (renderService *RenderService) GetMainRenderFile(postID uint) (string, *flock.Flock, error, error) {
	var filePath string

	// get post
	post, err := renderService.PostRepository.GetByID(postID)

	if err != nil {
		return filePath, nil, nil, fmt.Errorf("failed to find post with id %v: %w", postID, err)
	}

	// if render is pending return 202
	if post.RenderStatus == models.Pending {
		return filePath, nil, fmt.Errorf("render is still pending"), nil
	}

	// if render is failed return 404
	if post.RenderStatus == models.Failure {
		return filePath, nil, nil, fmt.Errorf("render has failed")
	}

	// lock directory
	// unlock upon error or after controller has read file
	lock, err := renderService.Filesystem.LockDirectory(postID)
	if err != nil {
		return filePath, nil, nil, fmt.Errorf("failed to aquire lock for directory %v: %w", postID, err)
	}

	// select repository of the post
	renderService.Filesystem.CheckoutDirectory(postID)

	// checkout master
	if err := renderService.Filesystem.CheckoutBranch("master"); err != nil {
		lock.Unlock()
		return filePath, nil, nil, fmt.Errorf("failed to find master: %w", err)
	}

	// verify render exists. if it doesn't set render status to failed
	fileName, err := renderService.Filesystem.RenderExists()
	if err != nil {
		post.RenderStatus = models.Failure
		_, _ = renderService.PostRepository.Update(post)
		lock.Unlock()

		return filePath, nil, nil, fmt.Errorf("render has failed: %w", err)
	}

	// set filepath to absolute path to render file
	filePath = filepath.Join(renderService.Filesystem.GetCurrentRenderDirPath(), fileName)

	return filePath, lock, nil, nil
}

func (renderService *RenderService) RenderPost(post *models.Post, lock *flock.Flock) {
	// defer unlocking repo
	defer lock.Unlock()

	// Checkout master
	if err := renderService.Filesystem.CheckoutBranch("master"); err != nil {
		renderService.FailPost(post)

		return
	}
	// Unzip saved file
	if err := renderService.Filesystem.Unzip(); err != nil {
		renderService.FailPost(post)

		return
	}

	// Validate project
	if valid := renderService.IsValidProject(); !valid {
		renderService.FailPost(post)

		return
	}

	// Install dependencies
	if err := renderService.InstallRenderDependencies(); err != nil {
		renderService.FailPost(post)

		return
	}

	// Set custom render config in yaml
	if err := renderService.SetProjectConfig(); err != nil {
		renderService.FailPost(post)

		return
	}

	// Render quarto project
	if err := renderService.RunRender(); err != nil {
		renderService.FailPost(post)

		return
	}

	// Verify that a render was produced in the form of a single file
	if _, err := renderService.Filesystem.RenderExists(); err != nil {
		renderService.FailPost(post)

		return
	}

	// Commit
	if err := renderService.Filesystem.CreateCommit(); err != nil {
		renderService.FailPost(post)

		return
	}

	// Update post render status
	post.RenderStatus = models.Success
	if _, err := renderService.PostRepository.Update(post); err != nil {
		renderService.FailPost(post)

		return
	}
}

func (renderService *RenderService) RenderBranch(branch *models.Branch, lock *flock.Flock) {
	// defer unlocking the repository
	defer lock.Unlock()

	// Checkout the branch
	if err := renderService.Filesystem.CheckoutBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		renderService.FailBranch(branch)
		_ = renderService.Filesystem.Reset()

		return
	}

	// Unzip saved file
	if err := renderService.Filesystem.Unzip(); err != nil {
		renderService.FailBranch(branch)
		_ = renderService.Filesystem.Reset()

		return
	}

	// Validate project
	if valid := renderService.IsValidProject(); !valid {
		renderService.FailBranch(branch)
		_ = renderService.Filesystem.Reset()

		return
	}

	// Install dependencies
	if err := renderService.InstallRenderDependencies(); err != nil {
		renderService.FailBranch(branch)
		_ = renderService.Filesystem.Reset()

		return
	}

	// Set custom render config in yaml
	if err := renderService.SetProjectConfig(); err != nil {
		renderService.FailBranch(branch)
		_ = renderService.Filesystem.Reset()

		return
	}

	// Render quarto project
	if err := renderService.RunRender(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	// Verify that a render was produced in the form of a single file
	if _, err := renderService.Filesystem.RenderExists(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	// Commit
	if err := renderService.Filesystem.CreateCommit(); err != nil {
		renderService.FailBranch(branch)

		return
	}

	// Update branch render status
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
			return fmt.Errorf("%s: %w", string(out), err)
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
			return fmt.Errorf("%s: %w", string(out), err)
		}

		// Install knitr
		cmd = exec.Command("Rscript", "-e", "renv::install('knitr')")
		cmd.Dir = renderService.Filesystem.GetCurrentQuartoDirPath()
		out, err = cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("%s: %w", string(out), err)
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
		return fmt.Errorf("failed to open yaml config file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, yamlObj)

	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml config file: %w", err)
	}

	yamlObj["format"] = map[string]interface{}{"html": map[string]interface{}{"page-layout": "custom"}}
	yamlFile, err = yaml.Marshal(yamlObj)

	if err != nil {
		return fmt.Errorf("failed to marshal yaml config file: %w", err)
	}

	var permMode fs.FileMode = 0o666
	err = os.WriteFile(configFilepath, yamlFile, permMode)

	if err != nil {
		return fmt.Errorf("failed to write yaml config file back: %w", err)
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
		return fmt.Errorf("quarto failed to render:\n%v\nerror: %w", out, err)
	}

	return nil
}

func (renderService *RenderService) FailBranch(branch *models.Branch) {
	branch.RenderStatus = models.Failure
	_, _ = renderService.BranchRepository.Update(branch)
}

func (renderService *RenderService) FailPost(post *models.Post) {
	post.RenderStatus = models.Failure
	_, _ = renderService.PostRepository.Update(post)
	_ = renderService.Filesystem.Reset()
}
