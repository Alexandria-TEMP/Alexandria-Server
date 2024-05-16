package filesystem_tests

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
	Filesystem = filesystem.InitFilesystem()
	cwd, _ = os.Getwd()

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
	file := filesystem.CreateMultipartFileHeader("./util/test.zip")

	Filesystem.SetCurrentVersion(1, 2)

	// Test saving fileheader
	err := Filesystem.SaveRepository(c, file)
	assert.Equal(t, nil, err)

	zipFilePath := filepath.Join(cwd, "vfs", "2", "1", "quarto_project.zip")
	_, err = os.Stat(zipFilePath)
	assert.Equal(t, nil, err)

	// Test unzipping succeeds and that contents are correct
	err = Filesystem.Unzip()
	assert.Equal(t, nil, err)

	projectFilePath := filepath.Join(cwd, "vfs", "2", "1", "quarto_project", "1234.txt")
	f, err := os.Open(projectFilePath)
	assert.Equal(t, nil, err)

	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := Readln(reader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "5678", line)

	// Test removing project directory
	err = Filesystem.RemoveProjectDirectory()
	assert.Equal(t, nil, err)

	projectDirPath := filepath.Join(cwd, "vfs", "2", "1", "quarto_project")
	_, err = os.Stat(projectDirPath)
	assert.NotEqual(t, nil, err)
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
