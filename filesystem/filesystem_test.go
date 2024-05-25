package filesystem

import (
	"bufio"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	CurrentFilesystem Filesystem
	cwdTest           string
)

func TestMain(m *testing.M) {
	cwdTest, _ = os.Getwd()

	CurrentFilesystem = *InitFilesystem()

	os.Exit(m.Run())
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwdTest, "vfs"))
}

func TestInitsystem(t *testing.T) {
	defer cleanup(t)

	CurrentFilesystem.SetCurrentVersion(1, 2)

	assert.Equal(t, filepath.Join(cwdTest, "vfs", "2", "1"), CurrentFilesystem.CurrentDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "2", "1", "quarto_project"), CurrentFilesystem.CurrentQuartoDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "2", "1", "render"), CurrentFilesystem.CurrentRenderDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "2", "1", "quarto_project.zip"), CurrentFilesystem.CurrentZipFilePath)
}

func TestFileHandling(t *testing.T) {
	defer cleanup(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	file, _ := CreateMultipartFileHeader("../utils/test_files/file_handling_test.zip")

	CurrentFilesystem.SetCurrentVersion(1, 2)

	// Test saving fileheader
	err := CurrentFilesystem.SaveRepository(c, file)
	assert.Nil(t, err)
	assert.True(t, FileExists(CurrentFilesystem.CurrentZipFilePath))

	// Test unzipping succeeds and that contents are correct
	err = CurrentFilesystem.Unzip()
	assert.Nil(t, err)

	projectFilePath := filepath.Join(CurrentFilesystem.CurrentQuartoDirPath, "1234.txt")
	f, err := os.Open(projectFilePath)
	assert.Nil(t, err)

	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := Readln(reader)
	assert.Nil(t, err)
	assert.Equal(t, "5678", line)

	// Test removing project directory
	err = CurrentFilesystem.RemoveProjectDirectory()
	assert.Nil(t, err)
	assert.False(t, FileExists(CurrentFilesystem.CurrentQuartoDirPath))
}

func TestGetRenderFileSuccess(t *testing.T) {
	CurrentFilesystem.CurrentDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "good_repository_setup")
	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(CurrentFilesystem.CurrentDirPath, "render")
	CurrentFilesystem.CurrentZipFilePath = filepath.Join(CurrentFilesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	file, err := CurrentFilesystem.GetRenderFile()
	assert.Nil(t, err)
	assert.Equal(t, []byte{53, 54, 55, 56}, file)
}

func TestGetRenderFileFailure1(t *testing.T) {
	CurrentFilesystem.CurrentDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "bad_repository_setup_1")
	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(CurrentFilesystem.CurrentDirPath, "render")
	CurrentFilesystem.CurrentZipFilePath = filepath.Join(CurrentFilesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, err := CurrentFilesystem.GetRenderFile()
	assert.NotNil(t, err)
}

func TestGetRenderFileFailure2(t *testing.T) {
	CurrentFilesystem.CurrentDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "bad_repository_setup_2")
	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(CurrentFilesystem.CurrentDirPath, "render")
	CurrentFilesystem.CurrentZipFilePath = filepath.Join(CurrentFilesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, err := CurrentFilesystem.GetRenderFile()
	assert.NotNil(t, err)
}

func TestGetRepositoryFileSuccess(t *testing.T) {
	CurrentFilesystem.CurrentDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "good_repository_setup")
	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(CurrentFilesystem.CurrentDirPath, "render")
	CurrentFilesystem.CurrentZipFilePath = filepath.Join(CurrentFilesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRepositoryFile
	fileForm, contentType, err := CurrentFilesystem.GetRepositoryFile()
	assert.Nil(t, err)
	assert.NotNil(t, fileForm)
	assert.Equal(t, "multipart/form-data; boundary=", contentType[:30])
}

func TestGetRepositoryFileFailure1(t *testing.T) {
	CurrentFilesystem.CurrentDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "bad_repository_setup_1")
	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(CurrentFilesystem.CurrentDirPath, "render")
	CurrentFilesystem.CurrentZipFilePath = filepath.Join(CurrentFilesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, _, err := CurrentFilesystem.GetRepositoryFile()
	assert.NotNil(t, err)
}

func TestGetRepositoryFileFailure2(t *testing.T) {
	CurrentFilesystem.CurrentDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "bad_repository_setup_2")
	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(CurrentFilesystem.CurrentDirPath, "render")
	CurrentFilesystem.CurrentZipFilePath = filepath.Join(CurrentFilesystem.CurrentDirPath, "quarto_project.zip")

	// Test GetRenderFile
	_, _, err := CurrentFilesystem.GetRepositoryFile()
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
