package filesystem

import (
	"bufio"
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
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

	CurrentFilesystem.SetCurrentVersion(1)

	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1"), CurrentFilesystem.CurrentDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1", "quarto_project"), CurrentFilesystem.CurrentQuartoDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1", "render"), CurrentFilesystem.CurrentRenderDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1", "quarto_project.zip"), CurrentFilesystem.CurrentZipFilePath)
}

func TestFileHandling(t *testing.T) {
	defer cleanup(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	file, _ := CreateMultipartFileHeader("../utils/test_files/file_handling_test.zip")

	CurrentFilesystem.SetCurrentVersion(1)

	// Test saving fileheader
	err := CurrentFilesystem.SaveRepository(c, file)
	assert.Nil(t, err)
	assert.True(t, utils.FileExists(CurrentFilesystem.CurrentZipFilePath))

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

	// Test removing version repository
	err = CurrentFilesystem.RemoveRepository()
	assert.Nil(t, err)
	assert.False(t, utils.FileExists(CurrentFilesystem.CurrentDirPath))
}

func TestGetFileTreeSuccess(t *testing.T) {
	CurrentFilesystem.CurrentQuartoDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "file_tree")

	files, err := CurrentFilesystem.GetFileTree()

	assert.Nil(t, err)

	assert.Equal(t, map[string]int64{".": -1, "child_dir": -1, "child_dir/test.txt": 0, "example.qmd": 0}, files)
}

func TestGetFileTreeFailure(t *testing.T) {
	CurrentFilesystem.CurrentQuartoDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "file_tree", "doesntexist")

	_, err := CurrentFilesystem.GetFileTree()

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

// CreateMultipartFileHeader is used for testing, to simulate an incoming request with a file
func CreateMultipartFileHeader(filePath string) (*multipart.FileHeader, error) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))

	if err != nil {
		return nil, err
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		return nil, err
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	maxMemoryBits := 20
	multipartForm, err := formReader.ReadForm(1 << maxMemoryBits)

	if err != nil {
		return nil, err
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, err
	}

	return files[0], nil
}