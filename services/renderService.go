package services

import (
	"fmt"
	"io/fs"
	"log"
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

//nolint:gocritic // we need all three types of return errors to be granular
func (renderService *RenderService) GetRenderFile(branchID uint) (filePath string, lock *flock.Flock, err202 error, err204 error, err404 error) {
	// get branch
	branch, err := renderService.BranchRepository.GetByID(branchID)

	if err != nil {
		return filePath, nil, nil, nil, fmt.Errorf("failed to find branch with id %v: %w", branchID, err)
	}

	// get project post
	projectPost, err := renderService.BranchService.GetBranchProjectPost(branch)

	if err != nil {
		return filePath, nil, nil, nil, fmt.Errorf("failed to find project post with id %v: %w", branch.ProjectPostID, err)
	}

	// if render is pending return 202
	if branch.RenderStatus == models.Pending {
		return filePath, nil, fmt.Errorf("render is still pending"), nil, nil
	}

	// if render is failed return 404
	if branch.RenderStatus == models.Failure {
		return filePath, nil, nil, fmt.Errorf("render has failed"), nil
	}

	// lock directory
	// unlock upon error or after controller has read file
	lock, err = renderService.Filesystem.LockDirectory(projectPost.PostID)
	if err != nil {
		return filePath, nil, nil, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", projectPost.PostID, err)
	}

	// select repository of the parent post
	directoryFilesystem := renderService.Filesystem.CheckoutDirectory(projectPost.PostID)

	// checkout specified branch
	if err := directoryFilesystem.CheckoutBranch(fmt.Sprintf("%v", branchID)); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return filePath, nil, nil, nil, fmt.Errorf("failed to find this git branch, with name %v: %w", branchID, err)
	}

	// verify render exists. if it doesn't set render status to failed
	fileName, err := directoryFilesystem.RenderExists()
	if err != nil {
		branch.RenderStatus = models.Failure
		_, _ = renderService.BranchRepository.Update(branch)

		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return filePath, nil, nil, nil, fmt.Errorf("render has failed: %w", err)
	}

	// set filepath to absolute path to render file
	filePath = filepath.Join(directoryFilesystem.GetCurrentRenderDirPath(), fileName)

	return filePath, lock, nil, nil, nil
}

//nolint:gocritic // we need all three types of return errors to be granular
func (renderService *RenderService) GetMainRenderFile(postID uint) (filePath string, lock *flock.Flock, err202 error, err204 error, err404 error) {
	// get post
	post, err := renderService.PostRepository.GetByID(postID)

	if err != nil {
		return filePath, nil, nil, nil, fmt.Errorf("failed to find post with id %v: %w", postID, err)
	}

	// if render is pending return 202
	if post.RenderStatus == models.Pending {
		return filePath, nil, fmt.Errorf("render is still pending"), nil, nil
	}

	// if render is failed return 404
	if post.RenderStatus == models.Failure {
		return filePath, nil, nil, fmt.Errorf("render has failed"), nil
	}

	// lock directory
	// unlock upon error or after controller has read file
	lock, err = renderService.Filesystem.LockDirectory(postID)
	if err != nil {
		return filePath, nil, nil, nil, fmt.Errorf("failed to acquire lock for directory %v: %w", postID, err)
	}

	// select repository of the post
	directoryFilesystem := renderService.Filesystem.CheckoutDirectory(postID)

	// checkout master
	if err := directoryFilesystem.CheckoutBranch("master"); err != nil {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return filePath, nil, nil, nil, fmt.Errorf("failed to find master: %w", err)
	}

	// verify render exists. if it doesn't set render status to failed
	fileName, err := directoryFilesystem.RenderExists()
	if err != nil {
		post.RenderStatus = models.Failure
		_, _ = renderService.PostRepository.Update(post)

		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}

		return filePath, nil, nil, nil, fmt.Errorf("render has failed: %w", err)
	}

	// set filepath to absolute path to render file
	filePath = filepath.Join(directoryFilesystem.GetCurrentRenderDirPath(), fileName)

	return filePath, lock, nil, nil, nil
}

func (renderService *RenderService) RenderPost(post *models.Post, lock *flock.Flock, directoryFilesystem filesystemInterfaces.Filesystem) {
	// defer unlocking repo
	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// Checkout master
	if err := directoryFilesystem.CheckoutBranch("master"); err != nil {
		log.Printf("POST RENDER ERROR: failed to checkout master: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}
	// Unzip saved file
	if err := directoryFilesystem.Unzip(); err != nil {
		log.Printf("POST RENDER ERROR: failed to unzip: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Validate project
	if valid := renderService.isValidProject(directoryFilesystem); !valid {
		log.Printf("POST RENDER ERROR: invalid project")
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Install dependencies
	if err := renderService.installRenderDependencies(directoryFilesystem); err != nil {
		log.Printf("POST RENDER ERROR: failed to install dependencies: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Set custom render config in yaml
	if err := renderService.setProjectConfig(directoryFilesystem); err != nil {
		log.Printf("POST RENDER ERROR: failed to set project config: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Render quarto project
	if err := renderService.runRender(directoryFilesystem); err != nil {
		log.Printf("POST RENDER ERROR: failed to run render: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Verify that a render was produced in the form of a single file
	if _, err := directoryFilesystem.RenderExists(); err != nil {
		log.Printf("POST RENDER ERROR: render does not exist: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Commit
	if err := directoryFilesystem.CreateCommit(); err != nil {
		log.Printf("POST RENDER ERROR: failed to create commit: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}

	// Update post render status
	post.RenderStatus = models.Success
	if _, err := renderService.PostRepository.Update(post); err != nil {
		log.Printf("POST RENDER ERROR: failed to update post render status: %s", err)
		renderService.failPost(post, directoryFilesystem)

		return
	}
}

func (renderService *RenderService) RenderBranch(branch *models.Branch, lock *flock.Flock, directoryFilesystem filesystemInterfaces.Filesystem) {
	// defer unlocking the repository
	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Printf("Failed to unlock %s", lock.Path())
		}
	}()

	// Checkout the branch
	if err := directoryFilesystem.CheckoutBranch(fmt.Sprintf("%v", branch.ID)); err != nil {
		renderService.failBranch(branch, directoryFilesystem)

		log.Printf("BRANCH RENDER ERROR: failed to checkout branch: %s", err)

		return
	}

	// Unzip saved file
	if err := directoryFilesystem.Unzip(); err != nil {
		renderService.failBranch(branch, directoryFilesystem)

		log.Printf("BRANCH RENDER ERROR: failed to unzip: %s", err)

		return
	}

	// Validate project
	if valid := renderService.isValidProject(directoryFilesystem); !valid {
		renderService.failBranch(branch, directoryFilesystem)

		log.Printf("BRANCH RENDER ERROR: invalid project")

		return
	}

	// Install dependencies
	if err := renderService.installRenderDependencies(directoryFilesystem); err != nil {
		renderService.failBranch(branch, directoryFilesystem)

		log.Printf("BRANCH RENDER ERROR: failed to install dependencies: %s", err)

		return
	}

	// Set custom render config in yaml
	if err := renderService.setProjectConfig(directoryFilesystem); err != nil {
		renderService.failBranch(branch, directoryFilesystem)

		log.Printf("BRANCH RENDER ERROR: failed to set project config: %s", err)

		return
	}

	// Render quarto project
	if err := renderService.runRender(directoryFilesystem); err != nil {
		renderService.failBranch(branch, directoryFilesystem)
		log.Printf("BRANCH RENDER ERROR: failed to run render: %s", err)

		return
	}

	// Verify that a render was produced in the form of a single file
	if _, err := directoryFilesystem.RenderExists(); err != nil {
		renderService.failBranch(branch, directoryFilesystem)
		log.Printf("BRANCH RENDER ERROR: render does not exist: %s", err)

		return
	}

	// Commit
	if err := directoryFilesystem.CreateCommit(); err != nil {
		renderService.failBranch(branch, directoryFilesystem)
		log.Printf("BRANCH RENDER ERROR: failed to create commit: %s", err)

		return
	}

	// Update branch render status
	branch.RenderStatus = models.Success
	if _, err := renderService.BranchRepository.Update(branch); err != nil {
		renderService.failBranch(branch, directoryFilesystem)
		log.Printf("BRANCH RENDER ERROR: failed to update branch render status: %s", err)

		return
	}
}

// IsValidProject validates that the files are a valid default quarto project
// They mus have a _quarto.yml or _quarto.yaml file.
// They must be of default quarto project type.
func (renderService *RenderService) isValidProject(directoryFilesystem filesystemInterfaces.Filesystem) bool {
	// If there is no yml file to cofigure the project it is invalid
	ymlPath := filepath.Join(directoryFilesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
	yamlPath := filepath.Join(directoryFilesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")

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
func (renderService *RenderService) installRenderDependencies(directoryFilesystem filesystemInterfaces.Filesystem) error {
	// Check if renv.lock exists and if so get dependencies
	rLockPath := filepath.Join(directoryFilesystem.GetCurrentQuartoDirPath(), "renv.lock")
	if utils.FileExists(rLockPath) {
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = directoryFilesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("%s: %w", string(out), err)
		}
	}

	// Check if any renv exists.
	renvActivatePath := filepath.Join(directoryFilesystem.GetCurrentQuartoDirPath(), "renv", "activate.R")
	if utils.FileExists(renvActivatePath) {
		// Install rmarkdown
		cmd := exec.Command("Rscript", "-e", "renv::install('rmarkdown')")
		cmd.Dir = directoryFilesystem.GetCurrentQuartoDirPath()
		out, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("%s: %w", string(out), err)
		}

		// Install knitr
		cmd = exec.Command("Rscript", "-e", "renv::install('knitr')")
		cmd.Dir = directoryFilesystem.GetCurrentQuartoDirPath()
		out, err = cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("%s: %w", string(out), err)
		}
	}

	return nil
}

func (renderService *RenderService) setProjectConfig(directoryFilesystem filesystemInterfaces.Filesystem) error {
	// Find config file
	yamlFilepath := filepath.Join(directoryFilesystem.GetCurrentQuartoDirPath(), "_quarto.yaml")
	ymlFilepath := filepath.Join(directoryFilesystem.GetCurrentQuartoDirPath(), "_quarto.yml")
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
func (renderService *RenderService) runRender(directoryFilesystem filesystemInterfaces.Filesystem) error {
	// Run render command
	cmd := exec.Command("quarto", "render", directoryFilesystem.GetCurrentQuartoDirPath(),
		"--output-dir", directoryFilesystem.GetCurrentRenderDirPath(),
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

func (renderService *RenderService) failBranch(branch *models.Branch, directoryFilesystem filesystemInterfaces.Filesystem) {
	branch.RenderStatus = models.Failure
	_, _ = renderService.BranchRepository.Update(branch)
	_ = directoryFilesystem.Reset()
}

func (renderService *RenderService) failPost(post *models.Post, directoryFilesystem filesystemInterfaces.Filesystem) {
	post.RenderStatus = models.Failure
	_, _ = renderService.PostRepository.Update(post)
	_ = directoryFilesystem.Reset()
}
