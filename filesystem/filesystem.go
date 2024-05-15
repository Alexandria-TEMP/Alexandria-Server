package filesystem

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Filesystem struct {
	rootPath             string
	zipName              string
	quartoDirectoryName  string
	CurrentDirPath       string
	CurrentQuartoDirPath string
	CurrentZipFilePath   string
	CurrentRenderDirPath string
	renderFormat         string
}

var cwd, _ = os.Getwd()
var DefaultRootPath = filepath.Clean(filepath.Join(cwd, "vfs"))
var DefaultZipName = "quarto_project.zip"
var DefaultQuartoDirectoryName = "quarto_project"
var DefaultRenderFormat = "html"

func InitFilesystem() Filesystem {
	filesystem := Filesystem{
		rootPath:            DefaultRootPath,
		zipName:             DefaultZipName,
		quartoDirectoryName: DefaultQuartoDirectoryName,
		renderFormat:        DefaultRenderFormat,
	}

	err := os.MkdirAll(filesystem.rootPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return filesystem
}

func (filesystem *Filesystem) SetCurrentVersion(versionID, postID uint) {
	filesystem.CurrentDirPath = filepath.Join(filesystem.rootPath, strconv.FormatUint(uint64(postID), 10), strconv.FormatUint(uint64(versionID), 10))
	filesystem.CurrentQuartoDirPath = filepath.Join(filesystem.CurrentDirPath, filesystem.quartoDirectoryName)
	filesystem.CurrentZipFilePath = filepath.Join(filesystem.CurrentDirPath, filesystem.zipName)
	filesystem.CurrentRenderDirPath = filepath.Join(filesystem.CurrentDirPath, "render")
}

// SaveRepository saves a zip file to a ./vfs/{postID}/{versionID} in the filesystem and return the path to the directory.
func (filesystem *Filesystem) SaveRepository(c *gin.Context, file *multipart.FileHeader) error {
	// Save zip file
	err := c.SaveUploadedFile(file, filesystem.CurrentZipFilePath)

	if err != nil {
		return err
	}

	return nil
}

// Unzip will unzip the quarto_project.zip file, if present, of any post version.
// Errors if there is no such file or it can't unzip it.
func (filesystem *Filesystem) Unzip() error {
	archive, err := zip.OpenReader(filesystem.CurrentZipFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(filesystem.CurrentQuartoDirPath, f.Name)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)

			if err != nil {
				return err
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func (filesystem *Filesystem) RenderProject() error {
	filesystem.InstallRenderDependencies()

	cmd := exec.Command("quarto", "render", filesystem.CurrentQuartoDirPath, "--output-dir", filesystem.CurrentRenderDirPath, "--to", "html", "--no-cache")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Printf("%s %s output:\n%s\n", cmd.Path, cmd.Args, out)

	return nil
}

// InstallRenderDependencies first checks if a renv.lock file is present and if so gets all dependencies.
// Next it checks for the
func (filesystem *Filesystem) InstallRenderDependencies() error {
	// Check if renv exists
	rLockPath := filepath.Join(filesystem.CurrentQuartoDirPath, "renv.lock")
	if _, err := os.Stat(rLockPath); err == nil {
		// Install all existing dependencies from renv.lock
		cmd := exec.Command("Rscript", "-e", "renv::restore()")
		cmd.Dir = filesystem.CurrentQuartoDirPath
		err := cmd.Run()

		if err != nil {
			return err
		}
	}

	// Install rmarkdown
	cmd := exec.Command("Rscript", "-e", "renv::install('rmarkdown')")
	cmd.Dir = filesystem.CurrentQuartoDirPath
	err := cmd.Run()

	if err != nil {
		return err
	}

	// Install knitr
	cmd = exec.Command("Rscript", "-e", "renv::install('knitr')")
	cmd.Dir = filesystem.CurrentQuartoDirPath
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func (filesystem *Filesystem) RemoveProjectDirectory() error {
	err := os.RemoveAll(filesystem.CurrentQuartoDirPath)

	if err != nil {
		return err
	}

	err = os.Remove(filesystem.CurrentQuartoDirPath)

	if err != nil {
		return err
	}

	return nil
}
