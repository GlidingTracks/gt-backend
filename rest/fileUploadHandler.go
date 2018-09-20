package rest

import (
	"errors"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Used in debugging. TODO remove before prod
const fileNameFUH = "fileUploadHandler.go"

// FileUploadHandler holds the context and routes for this handler.
type FileUploadHandler struct {
	Ctx Context
	UploadFilePage string
}

// Bind sets up the routes to the mux router.
func (fuh FileUploadHandler) Bind(r *mux.Router) {
	r.HandleFunc("/upload", uploadFilePage).Methods("POST")
}

// Upload and save a file to the filesystem
func uploadFilePage(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	if uid == "" {
		gtbackend.DebugLog(fileNameFUH, "uploadFilePage", errors.New(constant.ErrorNoUIDProvided))

		http.Error(w, constant.ErrorNoUIDProvided, http.StatusBadRequest)
		return
	}

	httpCode, err := ProcessUploadRequest(r, uid)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), httpCode)
	}
}

// ProcessUploadRequest - Actual processing of the file upload
// Inspiration: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
func ProcessUploadRequest(r *http.Request, uid string) (int, error) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		return http.StatusBadRequest, err
	}

	defer file.Close()

	err = checkFileContentType(file, handler)
	if err != nil {
		return http.StatusUnsupportedMediaType, err
	}

	f, err := saveFileToFileSystem(uid, handler)
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		return http.StatusBadRequest, err
	}
	defer f.Close()

	io.Copy(f, file)

	return 0, nil
}

// Check whether or not a file is of type IGC
// https://golang.org/pkg/net/http/#DetectContentType
func checkFileContentType(file multipart.File, handler *multipart.FileHeader) error {
	buff := make([]byte, 512)

	if _, err := file.Read(buff); err != nil {
		gtbackend.DebugLog(fileNameFUH, "checkFileContentType", err)

		return err
	}

	content := http.DetectContentType(buff)

	if !strings.Contains(handler.Filename, "." + constant.IGCExtension) || !strings.Contains(content, constant.TextPlain) {
		return errors.New(constant.ErrorInvalidContentType)
	}

	return nil
}

// Save the uploaded file in the filesystem. Path: .Records/{uId}/
func saveFileToFileSystem(uid string, handler *multipart.FileHeader) (*os.File, error) {
	path := createFilePath(constant.LSRoot, uid)
	os.MkdirAll(path, os.ModePerm)

	// CleanedFileName
	cfn := cleanFilePath(handler.Filename)
	fileName := path + constant.Slash + cfn

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	return f, err
}

// Method for creating a new path OS independent
func createFilePath(args ...string) string {
	var path string

	for _, k := range args {
		path = filepath.Join(path, k)
	}

	return path
}

// If the user has supplied a filename with already existing filepath, clean it up
// and return only the filename
func cleanFilePath(filePath string) string {
	parts := strings.Split(filePath, constant.Slash)
	return parts[len(parts)-1]
}
