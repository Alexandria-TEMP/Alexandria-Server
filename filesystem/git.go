package filesystem

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CreateRepository a repo at the filesystems current dir
func (filesystem *Filesystem) CreateRepository() error {
	// get repository path
	directory := filesystem.GetCurrentDirPath()

	// clean directory
	os.RemoveAll(directory)

	// git init
	_, err := git.PlainInit(directory, false)

	if err != nil {
		return fmt.Errorf("failed git init: %w", err)
	}

	// set CurrentRepository to new repo
	filesystem.CurrentRepository, err = filesystem.CheckoutRepository()

	if err != nil {
		return fmt.Errorf("failed to open new repo: %w", err)
	}

	// create initial files
	cwd, _ := os.Getwd()
	var templateRepoPath string

	if strings.Split(cwd, "/")[len(strings.Split(cwd, "/"))-1] == "filesystem" {
		templateRepoPath = filepath.Join(cwd, "template_repo")
	} else {
		templateRepoPath = filepath.Join(cwd, "filesystem", "template_repo")
	}

	err = cp.Copy(templateRepoPath, filesystem.CurrentDirPath)

	if err != nil {
		return fmt.Errorf("failed to copy over default repository: %w", err)
	}

	// make initial commit
	if err := filesystem.CreateCommit(); err != nil {
		return fmt.Errorf("failed to make initial commit: %w", err)
	}

	return nil
}

// OpenRepository opens the repository at CurrentDirPath.
// Do this prior to any operations on a repo.
// If no repo has been initiated here this will error.
func (filesystem *Filesystem) CheckoutRepository() (*git.Repository, error) {
	r, err := git.PlainOpen(filesystem.CurrentDirPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return r, nil
}

// CreateBranch creates a new branch from the last commit on master with branchName as the name.
func (filesystem *Filesystem) CreateBranch(branchName string) error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get reference to commit we branch off of
	fromCommit, err := filesystem.GetLastCommit("master")

	if err != nil {
		return fmt.Errorf("failed to get last commit on master: %w", err)
	}

	// git checkout master
	// git branch <branchName>
	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), fromCommit.Hash())

	// save branch to .git
	if err = filesystem.CurrentRepository.Storer.SetReference(ref); err != nil {
		return fmt.Errorf("failed to save branch reference: %w", err)
	}

	return nil
}

func (filesystem *Filesystem) DeleteBranch(branchName string) error {
	// checkout master
	if err := filesystem.CheckoutBranch("master"); err != nil {
		return fmt.Errorf("failed to checkout branch master: %w", err)
	}

	// git branch -d <branchName>
	cmd := exec.Command("git", "branch", "-D", branchName)
	cmd.Dir = filesystem.CurrentDirPath

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete branch %s:\n%s", branchName, string(out))
	}

	return nil
}

func (filesystem *Filesystem) Merge(toMerge, mergeInto string) error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree: %w", err)
	}

	// checkout master to merge into it
	if err := filesystem.CheckoutBranch(mergeInto); err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", mergeInto, err)
	}

	// get last commit on <branchName>
	lastCommit, err := filesystem.GetLastCommit(toMerge)

	if err != nil {
		return fmt.Errorf("failed to get last commit on %s: %w", toMerge, err)
	}

	// git reset --hard <branchName>
	err = w.Reset(&git.ResetOptions{
		Commit: lastCommit.Hash(),
		Mode:   git.HardReset,
	})

	if err != nil {
		return fmt.Errorf("failed to reset master to %s: %w", toMerge, err)
	}

	// git clean --ffxd
	if err := w.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return fmt.Errorf("failed to clean master: %w", err)
	}

	return nil
}

func (filesystem *Filesystem) Reset() error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree: %w", err)
	}

	if err := w.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return fmt.Errorf("failed to reset current branch: %w", err)
	}

	return nil
}

// CreateCommit will stage and commit all changes to the current branch
func (filesystem *Filesystem) CreateCommit() error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree: %w", err)
	}

	// git add .
	if err = w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return fmt.Errorf("failed to stage all changes: %w", err)
	}

	// git commit -m "-"
	if _, err = w.Commit("-", &git.CommitOptions{AllowEmptyCommits: true}); err != nil {
		return fmt.Errorf("failed to commit staged changes: %w", err)
	}

	return nil
}

// CheckoutBranch switches to the latest commit of a branch and removes any untracked files.
func (filesystem *Filesystem) CheckoutBranch(branchName string) error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree: %w", err)
	}

	// git reset --hard
	if err := w.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return fmt.Errorf("failed to reset branch %s: %w", branchName, err)
	}

	// git checkout <branchName>
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  true,
	}

	if err := w.Checkout(&branchCoOpts); err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", branchName, err)
	}

	// git reset --hard
	if err := w.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return fmt.Errorf("failed to reset branch %s: %w", branchName, err)
	}

	// git clean --ffxd
	if err := w.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return fmt.Errorf("failed to clean branch %s: %w", branchName, err)
	}

	return nil
}

// GetMasterRef gets the hash for the last commit on master.
func (filesystem *Filesystem) GetLastCommit(branchName string) (*plumbing.Reference, error) {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return nil, fmt.Errorf("no repository is currently checked out")
	}

	ref, err := filesystem.CurrentRepository.Reference(plumbing.NewBranchReferenceName(branchName), true)

	if err != nil || ref.Type() == plumbing.InvalidReference {
		return nil, fmt.Errorf("failed to get master ref: %w", err)
	}

	return ref, nil
}

func (filesystem *Filesystem) CleanDir() error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree: %w", err)
	}

	// git add . (add all files, in order to track them)
	if err = w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return fmt.Errorf("failed to stage all changes: %w", err)
	}

	// git rm -rf . (then remove all tracked files)
	cmd := exec.Command("git", "rm", "-rf", ".")
	cmd.Dir = filesystem.CurrentDirPath

	if out, err := cmd.CombinedOutput(); err != nil && string(out) != "fatal: pathspec '.' did not match any files\n" {
		return fmt.Errorf("failed to remove currently staged files:\n%s", string(out))
	}

	return nil
}
