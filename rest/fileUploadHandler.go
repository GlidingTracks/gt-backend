package rest

import (
	"context"
	"errors"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"mime/multipart"
	"net/http"
	"strings"
)

// fileNameFUH - Used in debugging. TODO remove before prod
const fileNameFUH = "fileUploadHandler.go"
/*
// FileUploadHandler holds the context and routes for this handler.
type FileUploadHandler struct {
	Ctx            Context
	UploadFilePage string
}

// Bind sets up the routes to the mux router.
func (fuh FileUploadHandler) Bind(r *mux.Router) {
	r.HandleFunc(fuh.UploadFilePage, uploadFilePage).Methods("POST")
}

// uploadFilePage - Upload and save a file to the filesystem
func uploadFilePage(w http.ResponseWriter, r *http.Request) {
	httpCode, _, err := ProcessUploadRequest(r)
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "uploadFilePage", err)

		http.Error(w, err.Error(), httpCode)
	}
}
*/

// ProcessUploadRequest - Actual processing of the file upload
// Inspiration: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
func ProcessUploadRequest(app *firebase.App, r *http.Request) (httpCode int, payload models.FilePayload, err error) {
	uid := getUID(r)
	if uid == "" {
		gtbackend.DebugLog(fileNameFUH, "uploadFilePage", errors.New(constant.ErrorNoUIDProvided))

		err = errors.New(constant.ErrorNoUIDProvided)
		httpCode = http.StatusBadRequest
		return
	}

	r.ParseMultipartForm(32 << 20)

	src, handler, err := r.FormFile("file")
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		httpCode = http.StatusBadRequest
		return
	}

	defer src.Close()

	err = checkFileContentType(src, handler)
	if err != nil {
		httpCode = http.StatusUnsupportedMediaType
		return
	}

	f, p, err := gtbackend.SaveFileToLocalStorage(uid, handler.Filename, src)
	if err != nil {
		gtbackend.DebugLog(fileNameFUH, "ProcessUploadRequest", err)

		httpCode = http.StatusBadRequest
		return
	}

	defer f.Close()

	httpCode = http.StatusOK

	payload = models.FilePayload{
		UID:  uid,
		Path: p,
	}

	isPrivate := r.FormValue("private")
	bp := gtbackend.GetBoolFromString(isPrivate)

	md, lines, err := uploadMetadataToFirestore(app, payload, bp)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "insertTrackRecordPage", err)

		httpCode = http.StatusBadRequest
	}

	err = uploadFileToFirebase(app, md, lines)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "insertTrackRecordPage", err)

		httpCode = http.StatusBadRequest
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

	// Reset seeker in file
	file.Seek(0, 0)

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

// uploadToFirebase saves a FilePayload struct to the DB.
func uploadMetadataToFirestore(app *firebase.App, record models.FilePayload, isPrivate bool) (md models.IgcMetadata, lines []string, err error) {
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		return
	}

	cr, _, err := client.Collection(constant.CollectionTracks).Add(ctx, record)
	if err != nil {
		return
	}

	parser := gtbackend.Parser{
		Path: record.Path,
	}

	pIGC, lines := parser.Parse()

	md = models.IgcMetadata{
		Privacy: isPrivate,
		Time:    gtbackend.GetUnixTime(),
		UID:     record.UID,
		Record:  pIGC,
		TrackID: cr.ID,
	}

	// TODO, maybe validate md somehow before pushing it to db
	_, _, err = client.Collection(constant.IgcMetadata).Add(ctx, md)
	if err != nil {
		return
	}

	return
}

func uploadFileToFirebase(app *firebase.App, md models.IgcMetadata, lines []string) (err error) {
	client, err := app.Storage(context.Background())
	if err != nil {
		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return
	}

	var fileParsed strings.Builder
	for i := 0; i < len(lines); i++ {
		fileParsed.WriteString(lines[i])
		fileParsed.WriteRune('\n')
	}

	wc := bucket.Object(md.TrackID).NewWriter(context.Background())
	wc.ContentType = "text/plain"
	_, err = wc.Write([]byte(fileParsed.String()))
	err = wc.Close()

	return
}
