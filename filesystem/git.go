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
	filesystem.CurrentRepository, err = filesystem.CheckoutRepository()

	if err != nil {
		return fmt.Errorf("failed to open new repo")
	}

	// make initial commit
	filesystem.CreateCommit()

	return nil
}

// OpenRepository opens the repository at CurrentDirPath.
// Do this prior to any operations on a repo.
// If no repo has been initiated here this will error.
func (filesystem *Filesystem) CheckoutRepository() (*git.Repository, error) {
	if r, err := git.PlainOpen(filesystem.CurrentDirPath); err != nil {
		return nil, fmt.Errorf("failed to open repository")
	} else {
		return r, nil
	}
}

// CreateBranch creates a new branch from the last commit on master with branchName as the name.
func (filesystem *Filesystem) CreateBranch(branchName string) error {
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

// Megre actually resets master to the last commit on the branch we are merging
func (filesystem *Filesystem) Merge(toMerge string, mergeInto string) error {
	// get worktree
	w, err := filesystem.CurrentRepository.Worktree()

	if err != nil {
		return fmt.Errorf("failed to open worktree")
	}

	// checkout master to merge into it
	filesystem.CheckoutBranch(mergeInto)

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
	w.Clean(&git.CleanOptions{Dir: true})

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
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  false,
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
func (filesystem *Filesystem) GetLastCommit(branchName string) (*plumbing.Reference, error) {
	ref, err := filesystem.CurrentRepository.Reference(plumbing.NewBranchReferenceName(branchName), true)

	if err != nil || ref.Type() == plumbing.InvalidReference {
		return nil, fmt.Errorf("failed to get master ref")
	}

	return ref, nil
}
