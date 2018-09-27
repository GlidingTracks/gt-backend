package rest

import (
	"errors"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// fileNameFUH - Used in debugging. TODO remove before prod
const fileNameFUH = "fileUploadHandler.go"

// FileUploadHandler holds the context and routes for this handler.
type FileUploadHandler struct {
	Ctx            Context
	UploadFilePage string
}

// Bind sets up the routes to the mux router.
func (fuh FileUploadHandler) Bind(r *mux.Router) {
	r.HandleFunc("/upload", uploadFilePage).Methods("POST")
}

// uploadFilePage - Upload and save a file to the filesystem
func uploadFilePage(w http.ResponseWriter, r *http.Request) {
	httpCode, _, err := ProcessUploadRequest(r)
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "uploadFilePage", err)

		http.Error(w, err.Error(), httpCode)
	}
}

// ProcessUploadRequest - Actual processing of the file upload
// Inspiration: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
func ProcessUploadRequest(r *http.Request) (httpCode int, payload models.FilePayload, err error) {
	uid := getUID(r)
	if uid == "" {
		gtbackend.DebugLog(fileNameFUH, "uploadFilePage", errors.New(constant.ErrorNoUIDProvided))

		err = errors.New(constant.ErrorNoUIDProvided)
		httpCode = http.StatusBadRequest
		return
	}

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

	f, p, err := gtbackend.SaveFileToLocalStorage(uid, handler.Filename)
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		httpCode = http.StatusBadRequest
		return
	}
	defer f.Close()

	io.Copy(f, file)

	httpCode = http.StatusOK

	payload = models.FilePayload{
		UID:  uid,
		Path: p,
	}

	return
}

// checkFileContentType - Check whether or not a file is of type IGC
// https://golang.org/pkg/net/http/#DetectContentType
func checkFileContentType(file multipart.File, handler *multipart.FileHeader) (err error) {
	buff := make([]byte, 512)

	if _, err = file.Read(buff); err != nil {
		gtbackend.DebugLog(fileNameFUH, "checkFileContentType", err)

		return
	}

	content := http.DetectContentType(buff)

	if !strings.Contains(handler.Filename, "."+constant.IGCExtension) || !strings.Contains(content, constant.TextPlain) {
		err = errors.New(constant.ErrorInvalidContentType)
		return
	}

	return
}

// getUID retrieves the "uid" field from a multipart/form-data request.
func getUID(r *http.Request) (uid string) {
	uid = r.FormValue("uid")
	return
}
