package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func CreateMultipartFileHeader(filePath string) *multipart.FileHeader {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		log.Fatal(err)
		return nil
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		log.Fatal("multipart file not exists")
		return nil
	}

	return files[0]
}

func CreateMultipartFile(filePath string) (io.Reader, string) {
	body := new(bytes.Buffer)

	mwriter := multipart.NewWriter(body)
	defer mwriter.Close()

	w, _ := mwriter.CreateFormFile("file", filePath)

	in, _ := os.Open(filePath)
	defer in.Close()

	_, _ = io.Copy(w, in)

	return body, mwriter.FormDataContentType()
}

func Unzip(zipFilePath, outputPath string) error {
	archive, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(outputPath, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
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
