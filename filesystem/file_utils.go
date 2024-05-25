package filesystem

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

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

// CreateMultipartFile bundles a file into an object go can interact with to return it as a response.
// Returns file, content-type, and error
func CreateMultipartFile(filePath string) (io.Reader, string, error) {
	// create a buffer to hold the file in memory
	body := new(bytes.Buffer)

	mwriter := multipart.NewWriter(body)
	defer mwriter.Close()

	w, err := mwriter.CreateFormFile("file", filePath)

	if err != nil {
		return body, "", err
	}

	in, err := os.Open(filePath)

	if err != nil {
		return body, "", err
	}

	defer in.Close()

	_, err = io.Copy(w, in)

	if err != nil {
		return body, "", err
	}

	return body, mwriter.FormDataContentType(), nil
}

// CreateByteSliceFile converts file at filpath into blob
func CreateByteSliceFile(filepath string) ([]byte, error) {
	file, err := os.ReadFile(filepath) //read the content of file
	if err != nil {
		return nil, err
	}

	return file, nil
}

// FileExists checks that a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

// FileContains checks that a file contains a string
func FileContains(filePath, match string) bool {
	if !FileExists(filePath) {
		return false
	}

	text, _ := os.ReadFile(filePath)

	return strings.Contains(string(text), match)
}
