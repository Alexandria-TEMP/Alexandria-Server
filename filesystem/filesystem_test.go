package filesystem

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/utils"
)

var (
	CurrentFilesystem Filesystem
	cwdTest           string
)

func TestMain(m *testing.M) {
	cwdTest, _ = os.Getwd()

	CurrentFilesystem = *NewFilesystem()

	os.Exit(m.Run())
}

func cleanup(t *testing.T) {
	t.Helper()

	os.RemoveAll(filepath.Join(cwdTest, "vfs"))
}

func TestInitsystem(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	CurrentFilesystem.CheckoutDirectory(1)

	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1"), CurrentFilesystem.CurrentDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1", "quarto_project"), CurrentFilesystem.CurrentQuartoDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1", "render"), CurrentFilesystem.CurrentRenderDirPath)
	assert.Equal(t, filepath.Join(cwdTest, "vfs", "1", "quarto_project.zip"), CurrentFilesystem.CurrentZipFilePath)
	assert.Nil(t, CurrentFilesystem.CurrentRepository)
}

func TestGit(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	// Set current dir
	CurrentFilesystem.CheckoutDirectory(1)

	// Create repo
	assert.Nil(t, CurrentFilesystem.CreateRepository())

	// Create and checkout branch 1 from main
	assert.Nil(t, CurrentFilesystem.CreateBranch("1"))
	assert.Nil(t, CurrentFilesystem.CheckoutBranch("1"))

	// Add new file and commit
	helloFilePath := filepath.Join(CurrentFilesystem.GetCurrentDirPath(), "hello.txt")
	assert.Nil(t, os.WriteFile(helloFilePath, []byte("world"), fs.ModePerm))
	assert.Nil(t, CurrentFilesystem.CreateCommit())

	// Check file contents
	contents, _ := os.ReadFile(helloFilePath)
	assert.Equal(t, "world", string(contents))

	// Merge 1 into master and verify
	assert.Nil(t, CurrentFilesystem.Merge("1", "master"))

	contents, _ = os.ReadFile(helloFilePath)
	assert.Equal(t, "world", string(contents))

	// Create branch 2 and 3 from main at same time
	assert.Nil(t, CurrentFilesystem.CreateBranch("2"))
	assert.Nil(t, CurrentFilesystem.CreateBranch("3"))

	// Checkout branch 2, edit hello.txt, commit, merge, and verify changes
	assert.Nil(t, CurrentFilesystem.CheckoutBranch("2"))
	assert.Nil(t, os.WriteFile(helloFilePath, []byte("alexandria"), fs.ModePerm))
	assert.Nil(t, CurrentFilesystem.CreateCommit())
	assert.Nil(t, CurrentFilesystem.Merge("2", "master"))

	contents, _ = os.ReadFile(helloFilePath)
	assert.Equal(t, "alexandria", string(contents))

	// Checkout branch 3, delete hello.txt, add "README.md", commit, merge, and verify changes
	readmeFilePath := filepath.Join(CurrentFilesystem.GetCurrentDirPath(), "README.md")
	assert.Nil(t, CurrentFilesystem.CheckoutBranch("3"))
	assert.Nil(t, CurrentFilesystem.CleanDir())
	assert.Nil(t, os.WriteFile(readmeFilePath, []byte("welcome"), fs.ModePerm))
	assert.Nil(t, CurrentFilesystem.CreateCommit())
	assert.Nil(t, CurrentFilesystem.Merge("3", "master"))
	assert.False(t, utils.FileExists(helloFilePath))

	// Delete branch 3
	assert.Nil(t, CurrentFilesystem.DeleteBranch("2"))
	assert.NotNil(t, CurrentFilesystem.CheckoutBranch("2"))

	// Get last commit on master
	ref, err := CurrentFilesystem.GetLastCommit("master")
	assert.Nil(t, err)

	// Add files, reset before committing, and verify reset worked
	mistakeFilePath := filepath.Join(CurrentFilesystem.GetCurrentDirPath(), "oops")
	assert.Nil(t, os.WriteFile(mistakeFilePath, []byte("whoopsies"), fs.ModePerm))
	assert.Nil(t, CurrentFilesystem.Reset())

	ref2, err := CurrentFilesystem.GetLastCommit("master")
	assert.Nil(t, err)
	assert.Equal(t, ref, ref2)
	assert.False(t, utils.FileExists(mistakeFilePath))

	// Delete repo and verify it worked
	assert.Nil(t, CurrentFilesystem.DeleteRepository())

	_, err = CurrentFilesystem.CheckoutRepository()
	assert.NotNil(t, err)
}

func TestGitOperationsWithoutRepo(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	// Create branch without repo
	assert.NotNil(t, CurrentFilesystem.CreateBranch("master"))

	// Checkout branch without repo
	assert.NotNil(t, CurrentFilesystem.CheckoutBranch("master"))

	// Merge without repo
	assert.NotNil(t, CurrentFilesystem.Merge("2", "master"))

	// Reset without repo
	assert.NotNil(t, CurrentFilesystem.Reset())

	// Create commit without repo
	assert.NotNil(t, CurrentFilesystem.CreateCommit())

	// Get last commit without repo
	_, err := CurrentFilesystem.GetLastCommit("master")
	assert.NotNil(t, err)
}

func TestCheckoutNonexistantBranch(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	// Set current dir
	CurrentFilesystem.CheckoutDirectory(99)

	// Create repo
	assert.Nil(t, CurrentFilesystem.CreateRepository())

	// Checkout branch that doesnt exist
	assert.NotNil(t, CurrentFilesystem.CheckoutBranch("badbranch"))
}

func TestGitOperationsOnBareRepo(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	// get repository path
	directory := filepath.Join(cwdTest, "vfs", "1")

	// git init at
	_, err := git.PlainInit(directory, true)
	assert.Nil(t, err)

	// checkout bare repo
	CurrentFilesystem.CheckoutDirectory(1)

	// Create branch with bare repo
	assert.NotNil(t, CurrentFilesystem.CreateBranch("master"))

	// Checkout branch with bare repo
	assert.NotNil(t, CurrentFilesystem.CheckoutBranch("master"))

	// Merge with bare repo
	assert.NotNil(t, CurrentFilesystem.Merge("2", "master"))

	// Reset with bare repo
	assert.NotNil(t, CurrentFilesystem.Reset())

	// Create commit with bare repo
	assert.NotNil(t, CurrentFilesystem.CreateCommit())

	// Get last commit with bare repo
	_, err = CurrentFilesystem.GetLastCommit("master")
	assert.NotNil(t, err)
}

func TestFileHandling(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	file, _ := CreateMultipartFileHeader("../utils/test_files/file_handling_test.zip")

	CurrentFilesystem.CheckoutDirectory(1)

	// Test saving fileheader
	err := CurrentFilesystem.SaveZipFile(c, file)
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
	err = CurrentFilesystem.DeleteRepository()
	assert.Nil(t, err)
	assert.False(t, utils.FileExists(CurrentFilesystem.CurrentDirPath))
}

func TestUnzipDoesntExist(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	file, _ := CreateMultipartFileHeader("../utils/test_files/bad_zip")

	CurrentFilesystem.CheckoutDirectory(1)

	// Test saving fileheader
	err := CurrentFilesystem.SaveZipFile(c, file)
	assert.Nil(t, err)
	assert.True(t, utils.FileExists(CurrentFilesystem.CurrentZipFilePath))

	// Test unzipping succeeds and that contents are correct
	err = CurrentFilesystem.Unzip()
	assert.NotNil(t, err)
}

func TestRenderExistsSuccess(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "good_repository_setup", "render")

	exists, name := CurrentFilesystem.RenderExists()
	assert.True(t, exists)
	assert.Equal(t, "1234.html", name)
}

func TestRenderExistsNoFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "good_repository_setup", "badpath")

	exists, _ := CurrentFilesystem.RenderExists()
	assert.False(t, exists)
}

func TestRenderExistsMultipleFiles(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "bad_repository_setup_1")

	exists, _ := CurrentFilesystem.RenderExists()
	assert.False(t, exists)
}

func TestRenderExistsMultipleNotHtml(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	CurrentFilesystem.CurrentRenderDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "bad_repository_setup_2")

	exists, _ := CurrentFilesystem.RenderExists()
	assert.False(t, exists)
}

func TestGetFileTreeSuccess(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	defer cleanup(t)

	CurrentFilesystem.CurrentQuartoDirPath = filepath.Join(cwdTest, "..", "utils", "test_files", "file_tree")

	files, err := CurrentFilesystem.GetFileTree()

	assert.Nil(t, err)

	assert.Equal(t, map[string]int64{".": -1, "child_dir": -1, "child_dir/test.txt": 0, "example.qmd": 0}, files)
}

func TestGetFileTreeFailure(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

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
