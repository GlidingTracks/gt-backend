package rest

import (
	"cloud.google.com/go/firestore"
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

// fileNameFUH - Used in debugging.
const fileNameFUH = "fileUploadHandler.go"

// ProcessUploadRequest - Actual processing of the file upload
// Inspiration: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
func ProcessUploadRequest(app *firebase.App, r *http.Request) (httpCode int, md models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFUH, "ProcessUploadRequest")
	uid := getUID(r)

	if uid == "" {
		err = errors.New(constant.ErrorNoUIDProvided)
		gtbackend.DebugLogErrNoMsg(log, err)

		httpCode = http.StatusBadRequest
		return
	}

	r.ParseMultipartForm(32 << 20)

	// Get file source
	src, handler, err := r.FormFile("file")
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		httpCode = http.StatusBadRequest
		return
	}
	defer src.Close()

	// Check file and parse it to workable string
	parsed, err := processFileContent(src, handler)
	if err != nil {
		httpCode = http.StatusUnsupportedMediaType
		return
	}

	parsed, err = AnalyzeIGC(parsed)

	if err != nil {
		httpCode = http.StatusNoContent
		return
	}

	// Deduce file privacy
	isPrivate := r.FormValue("private")
	bp := gtbackend.GetBoolFromString(isPrivate)

	// Upload metadata to Cloud Firestore Database
	md, err = uploadMetadataToFirestore(app, uid, parsed, bp)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		httpCode = http.StatusBadRequest
		return
	}

	// Upload file to Firebase Storage
	err = uploadFileToFirebase(app, md, parsed)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		httpCode = http.StatusBadRequest
		return
	}

	httpCode = http.StatusOK
	return
}

// checkFileContentType - Check whether or not a file is of type IGC
// https://golang.org/pkg/net/http/#DetectContentType
func processFileContent(file multipart.File, handler *multipart.FileHeader) (parsed string, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFUH, "processFileContent")
	if handler.Size > constant.MaxIgcFileSize {
		gtbackend.DebugLogErrMsg(log, err, "MaxIgcFileSize FAIL")

		return
	}

	buff := make([]byte, handler.Size)

	if _, err = file.Read(buff); err != nil {
		gtbackend.DebugLogErrMsg(log, err, "Buffer error")

		return
	}

	// Reset seeker in file
	file.Seek(0, 0)

	content := http.DetectContentType(buff)

	// Check file name for extension and contents to be the TextPlain constant
	if !strings.Contains(handler.Filename, "."+constant.IGCExtension) || !strings.Contains(content, constant.TextPlain) {
		err = errors.New(constant.ErrorInvalidContentType)
		gtbackend.DebugLogErrMsg(log, err, "Content error")

		return
	}

	// Parse the file as we already buffer it here
	var builder strings.Builder
	builder.Write(buff)
	parsed = builder.String()

	return
}

// getUID retrieves the "uid" field from a multipart/form-data request.
func getUID(r *http.Request) (uid string) {
	uid = r.Header.Get("uid")
	return
}

// uploadMetadataToFirestore saves a FilePayload struct to the DB.
func uploadMetadataToFirestore(app *firebase.App, uid string, parsed string, isPrivate bool) (md models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFUH, "uploadMetadataToFirestore")
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	parser := gtbackend.Parser{
		Parsed: parsed,
	}

	pIGC, _ := parser.Parse()

	// Prepare md with data, TrackID as placeholder for now
	md = models.IgcMetadata{
		Privacy: isPrivate,
		Time:    gtbackend.GetUnixTime(),
		Uid:     uid,
		Record:  pIGC,
		TrackID: "placeholder",
	}

	// Upload metadata to Firestore, get document record to update TrackID later
	cr, _, err := client.Collection(constant.IgcMetadata).Add(ctx, md)
	if err != nil {
		gtbackend.DebugLogErrMsg(log, err, "Add error")

		return
	}

	// Set TrackID to be the document ID of the metadata
	_, err = client.Collection(constant.IgcMetadata).Doc(cr.ID).Set(ctx, map[string]interface{}{
		"TrackID": cr.ID,
	}, firestore.MergeAll)
	md.TrackID = cr.ID
	if err != nil {
		gtbackend.DebugLogErrMsg(log, err, "Set TrackID error")

		return
	}

	return
}

// uploadFileToFirebase Uploads the file to Firebase Storage
func uploadFileToFirebase(app *firebase.App, md models.IgcMetadata, parsed string) (err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFUH, "uploadFileToFirebase")
	client, err := app.Storage(context.Background())
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	//Create new object in storage with TrackID ID, write parsed string to the file
	wc := bucket.Object(md.TrackID).NewWriter(context.Background())
	defer wc.Close()

	wc.ContentType = "text/plain"
	_, err = wc.Write([]byte(parsed))
	if err != nil {
		gtbackend.DebugLogErrMsg(log, err, "Write error")
		return
	}

	return
}
