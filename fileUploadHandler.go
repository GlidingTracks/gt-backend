package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"os"
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

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Error(err)
		return err, http.StatusBadRequest
	}
	defer f.Close()

	//io.Copy(f, file)

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

func saveFileToFileSystem() {
	// TODO save to Firebase
}