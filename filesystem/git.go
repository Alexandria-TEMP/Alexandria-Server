package filesystem

import (
	"fmt"
	"os"

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
		return fmt.Errorf("failed it git init")
	}

	// set CurrentRepository to new repo
	filesystem.CurrentRepository, err = filesystem.OpenRepository()

	if err != nil {
		return fmt.Errorf("failed to open new repo")
	}

	return nil
}

// OpenRepository opens the repository at CurrentDirPath.
// Do this prior to any operations on a repo.
// If no repo has been initiated here this will error.
func (filesystem *Filesystem) OpenRepository() (*git.Repository, error) {
	if r, err := git.PlainOpen(filesystem.CurrentDirPath); err != nil {
		return nil, fmt.Errorf("failed to open repository")
	} else {
		return r, nil
	}
}

// CreateBranch creates a new branch from the last commit on master with branchName as the name.
func (filesystem *Filesystem) CreateBranch(branchName string) error {
	// get reference to commit we branch off of
	fromRef, err := filesystem.GetMasterRef()

	if err != nil {
		return err
	}

	// git checkout master
	// git branch <branchName>
	toRef := plumbing.NewHashReference(plumbing.ReferenceName(branchName), fromRef.Hash())

	// save branch to .git
	if err = filesystem.CurrentRepository.Storer.SetReference(toRef); err != nil {
		return fmt.Errorf("failed to create new branch")
	}

	return nil
}

// CreateCommit will stage and commit all changes to the current branch
func (filesystem *Filesystem) CreateCommit() error {
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
	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree")
	}

	// git checkout <branchName>
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchName),
		Force:  true,
	}

	if err := w.Checkout(&branchCoOpts); err != nil {
		return fmt.Errorf("failed to checkout branch %s", branchName)
	}

	// git reset --hard
	w.Reset(&git.ResetOptions{Mode: git.HardReset})

	// git clean --ffxd
	w.Clean(&git.CleanOptions{Dir: true})

	return nil
}

// GetMasterRef gets the hash for the last commit on master.
func (filesystem *Filesystem) GetMasterRef() (*plumbing.Reference, error) {
	if ref, err := filesystem.CurrentRepository.Reference(plumbing.Master, true); err != nil || ref.Type() == plumbing.InvalidReference {
		return nil, fmt.Errorf("failed to get master ref")
	} else {
		return ref, nil
	}
}
