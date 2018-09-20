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
func ProcessUploadRequest(r *http.Request, uid string) (httpCode int, err error) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		httpCode = http.StatusBadRequest
		return
	}

	defer file.Close()

	err = checkFileContentType(file, handler)
	if err != nil {
		httpCode = http.StatusUnsupportedMediaType
		return
	}

	f, err := saveFileToFileSystem(uid, handler)
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		httpCode = http.StatusBadRequest
		return
	}
	defer f.Close()

	io.Copy(f, file)

	httpCode = http.StatusOK
	return
}

// Check whether or not a file is of type IGC
// https://golang.org/pkg/net/http/#DetectContentType
func checkFileContentType(file multipart.File, handler *multipart.FileHeader) (err error) {
	buff := make([]byte, 512)

	if _, err = file.Read(buff); err != nil {
		gtbackend.DebugLog(fileNameFUH, "checkFileContentType", err)

		return
	}

	content := http.DetectContentType(buff)

	if !strings.Contains(handler.Filename, "." + constant.IGCExtension) || !strings.Contains(content, constant.TextPlain) {
		err = errors.New(constant.ErrorInvalidContentType)
		return
	}

	return
}

// Save the uploaded file in the filesystem. Path: .Records/{uId}/
func saveFileToFileSystem(uid string, handler *multipart.FileHeader) (file *os.File, err error) {
	path := createFilePath(constant.LSRoot, uid)
	os.MkdirAll(path, os.ModePerm)

	// CleanedFileName
	cfn := cleanFilePath(handler.Filename)
	fileName := path + constant.Slash + cfn

	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)

	return
}

// Method for creating a new path OS independent
func createFilePath(args ...string) (path string) {
	for _, k := range args {
		path = filepath.Join(path, k)
	}

	return
}

// If the user has supplied a filename with already existing filepath, clean it up
// and return only the filename
func cleanFilePath(filePath string) (fileName string) {
	parts := strings.Split(filePath, constant.Slash)

	fileName = parts[len(parts)-1]
	return
}
