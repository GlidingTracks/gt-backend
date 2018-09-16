package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Actual processing of the file upload
// Inspiration: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
func ProcessUploadRequest(r* http.Request, uid string) (error, int) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		logrus.Error("Could not get file: ", err)
		return err, http.StatusBadRequest
	}

	defer file.Close()

	err = checkFileContentType(file, handler)
	if err != nil {
		return err, http.StatusUnsupportedMediaType
	}

	f, err := saveFileToFileSystem(uid, handler)
	if err != nil {
		logrus.Error(err)
		return err, http.StatusBadRequest
	}
	defer f.Close()

	io.Copy(f, file)

	return nil, 0
}

// Check whether or not a file is of type IGC
// https://golang.org/pkg/net/http/#DetectContentType
func checkFileContentType(file multipart.File, handler* multipart.FileHeader) (error){
	buff := make([]byte, 512)

	if _, err := file.Read(buff); err != nil {
		logrus.Error(err)
		return err
	}

	content := http.DetectContentType(buff)

	if !strings.Contains(handler.Filename, ".igc") || !strings.Contains(content, "text/plain") {
		return errors.New("invalid content-type")
	}

	return nil
}

// Save the uploaded file in the filesystem. Path: .Records/{uId}/
func saveFileToFileSystem(uid string, handler* multipart.FileHeader) (*os.File, error) {
	path := createFilePath("Records", uid)
	os.MkdirAll(path, os.ModePerm)

	// CleanedFileName
	cfn := cleanFilePath(handler.Filename)
	fileName := path + "/" + cfn

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	return f, err
}

// Method for creating a new path OS independent
func createFilePath(args ...string) (string) {
	var path string

	for _, k := range args {
		path = filepath.Join(path, k)
	}

	return path
}

// If the user has supplied a filename with already existing filepath, clean it up
// and return only the filename
func cleanFilePath(filePath string) (string) {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}