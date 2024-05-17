package filesystemtests

import (
	"bufio"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/filesystem"
)

var (
	Filesystem filesystem.Filesystem
	cwd        string
)

func TestMain(m *testing.M) {
	cwd, _ = os.Getwd()

	Filesystem = *filesystem.InitFilesystem()

	os.Exit(m.Run())
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwd, "vfs"))
}

func TestInitsystem(t *testing.T) {
	defer cleanup(t)

	Filesystem.SetCurrentVersion(1, 2)

	assert.Equal(t, filepath.Join(cwd, "vfs", "2", "1"), Filesystem.CurrentDirPath)
	assert.Equal(t, filepath.Join(cwd, "vfs", "2", "1", "quarto_project"), Filesystem.CurrentQuartoDirPath)
	assert.Equal(t, filepath.Join(cwd, "vfs", "2", "1", "render"), Filesystem.CurrentRenderDirPath)
	assert.Equal(t, filepath.Join(cwd, "vfs", "2", "1", "quarto_project.zip"), Filesystem.CurrentZipFilePath)
}

func TestFileHandling(t *testing.T) {
	defer cleanup(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	file, _ := filesystem.CreateMultipartFileHeader("../utils/file_handling_test.zip")

	Filesystem.SetCurrentVersion(1, 2)

	// Test saving fileheader
	err := Filesystem.SaveRepository(c, file)
	assert.Nil(t, err)
	assert.True(t, filesystem.FileExists(Filesystem.CurrentZipFilePath))

	// Test unzipping succeeds and that contents are correct
	err = Filesystem.Unzip()
	assert.Nil(t, err)

	projectFilePath := filepath.Join(Filesystem.CurrentQuartoDirPath, "1234.txt")
	f, err := os.Open(projectFilePath)
	assert.Nil(t, err)

	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := Readln(reader)
	assert.Nil(t, err)
	assert.Equal(t, "5678", line)

	// Test removing project directory
	err = Filesystem.RemoveProjectDirectory()
	assert.Nil(t, err)
	assert.False(t, filesystem.FileExists(Filesystem.CurrentQuartoDirPath))
}

func TestGetRenderFileSuccess(t *testing.T) {
	Filesystem.CurrentDirPath = filepath.Join(cwd, "..", "utils", "good_repository_setup")
	Filesystem.CurrentRenderDirPath = filepath.Join(Filesystem.CurrentDirPath, "render")
	Filesystem.CurrentZipFilePath = filepath.Join(Filesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	fileForm, contentType, err := Filesystem.GetRenderFile()
	assert.Nil(t, err)
	assert.NotNil(t, fileForm)
	assert.Equal(t, "multipart/form-data; boundary=", contentType[:30])
}

func TestGetRenderFileFailure1(t *testing.T) {
	Filesystem.CurrentDirPath = filepath.Join(cwd, "..", "utils", "bad_repository_setup_1")
	Filesystem.CurrentRenderDirPath = filepath.Join(Filesystem.CurrentDirPath, "render")
	Filesystem.CurrentZipFilePath = filepath.Join(Filesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, _, err := Filesystem.GetRenderFile()
	assert.NotNil(t, err)
}

func TestGetRenderFileFailure2(t *testing.T) {
	Filesystem.CurrentDirPath = filepath.Join(cwd, "..", "utils", "bad_repository_setup_2")
	Filesystem.CurrentRenderDirPath = filepath.Join(Filesystem.CurrentDirPath, "render")
	Filesystem.CurrentZipFilePath = filepath.Join(Filesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, _, err := Filesystem.GetRenderFile()
	assert.NotNil(t, err)
}

func TestGetRepositoryFileSuccess(t *testing.T) {
	Filesystem.CurrentDirPath = filepath.Join(cwd, "..", "utils", "good_repository_setup")
	Filesystem.CurrentRenderDirPath = filepath.Join(Filesystem.CurrentDirPath, "render")
	Filesystem.CurrentZipFilePath = filepath.Join(Filesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRepositoryFile
	fileForm, contentType, err := Filesystem.GetRepositoryFile()
	assert.Nil(t, err)
	assert.NotNil(t, fileForm)
	assert.Equal(t, "multipart/form-data; boundary=", contentType[:30])
}

func TestGetRepositoryFileFailure1(t *testing.T) {
	Filesystem.CurrentDirPath = filepath.Join(cwd, "..", "utils", "bad_repository_setup_1")
	Filesystem.CurrentRenderDirPath = filepath.Join(Filesystem.CurrentDirPath, "render")
	Filesystem.CurrentZipFilePath = filepath.Join(Filesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, _, err := Filesystem.GetRepositoryFile()
	assert.NotNil(t, err)
}

func TestGetRepositoryFileFailure2(t *testing.T) {
	Filesystem.CurrentDirPath = filepath.Join(cwd, "..", "utils", "bad_repository_setup_2")
	Filesystem.CurrentRenderDirPath = filepath.Join(Filesystem.CurrentDirPath, "render")
	Filesystem.CurrentZipFilePath = filepath.Join(Filesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, _, err := Filesystem.GetRepositoryFile()
	assert.NotNil(t, err)
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}
