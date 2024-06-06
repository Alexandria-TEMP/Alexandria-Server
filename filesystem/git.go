package filesystem

import (
	"fmt"
	"os"
	"os/exec"

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
		return fmt.Errorf("failed git init")
	}

	// set CurrentRepository to new repo
	filesystem.CurrentRepository, err = filesystem.CheckoutRepository()

	if err != nil {
		return fmt.Errorf("failed to open new repo")
	}

	// make initial commit
	if err := filesystem.CreateCommit(); err != nil {
		return fmt.Errorf("failed to make initial commit")
	}

	return nil
}

// OpenRepository opens the repository at CurrentDirPath.
// Do this prior to any operations on a repo.
// If no repo has been initiated here this will error.
func (filesystem *Filesystem) CheckoutRepository() (*git.Repository, error) {
	r, err := git.PlainOpen(filesystem.CurrentDirPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open repository")
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
		return err
	}

	// git checkout master
	// git branch <branchName>
	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), fromCommit.Hash())

	// save branch to .git
	if err = filesystem.CurrentRepository.Storer.SetReference(ref); err != nil {
		return fmt.Errorf("failed to create new branch")
	}

	return nil
}

func (filesystem *Filesystem) DeleteBranch(branchName string) error {
	// checkout master
	if err := filesystem.CheckoutBranch("master"); err != nil {
		return fmt.Errorf("failed to checkout branch master")
	}

	// git branch -d <branchName>
	cmd := exec.Command("git", "branch", "-D", branchName)
	cmd.Dir = filesystem.CurrentDirPath

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete branch %v with message:\n %v", branchName, out)
	}

	return nil
}

// Merge actually resets master to the last commit on the branch we are merging
func (filesystem *Filesystem) Merge(toMerge, mergeInto string) error {
	// check if we have a repo open
	if filesystem.CurrentRepository == nil {
		return fmt.Errorf("no repository is currently checked out")
	}

	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree")
	}

	// checkout master to merge into it
	if err := filesystem.CheckoutBranch(mergeInto); err != nil {
		return fmt.Errorf("failed to checkout branch %v", mergeInto)
	}

	// get last commit on <branchName>
	lastCommit, err := filesystem.GetLastCommit(toMerge)

	if err != nil {
		return fmt.Errorf("failed to fetch last commit on %s", toMerge)
	}

	// git reset --hard <branchName>
	err = w.Reset(&git.ResetOptions{
		Commit: lastCommit.Hash(),
		Mode:   git.HardReset,
	})

	if err != nil {
		return fmt.Errorf("failed to reset master to %s", toMerge)
	}

	// git clean --ffxd
	if err := w.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return fmt.Errorf("failed to clean master")
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
		return fmt.Errorf("failed to open worktree")
	}

	if err := w.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return fmt.Errorf("failed to reset current branch")
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
		return fmt.Errorf("failed to open worktree")
	}

	// git add .
	if err = w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return fmt.Errorf("failed to stage all changes")
	}

	// git commit -m "-"
	if _, err = w.Commit("-", &git.CommitOptions{AllowEmptyCommits: true}); err != nil {
		return fmt.Errorf("failed to commit stages changes")
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
		return fmt.Errorf("failed to open worktree")
	}

	// git checkout <branchName>
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  false,
	}

	if err := w.Checkout(&branchCoOpts); err != nil {
		return fmt.Errorf("failed to checkout branch %s", branchName)
	}

	// git reset --hard
	if err := w.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return fmt.Errorf("failed to reset branch %s", branchName)
	}

	// git clean --ffxd
	if err := w.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return fmt.Errorf("failed to clean branch %s", branchName)
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
		return nil, fmt.Errorf("failed to get master ref")
	}

	return ref, nil
}

func (filesystem *Filesystem) CleanDir() error {
	// git rm -rf .
	cmd := exec.Command("git", "rm", "-rf", ".")
	cmd.Dir = filesystem.CurrentDirPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove all files from index")
	}

	return nil
}
