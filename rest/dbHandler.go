package rest

import (
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/gorilla/mux"
	"net/http"
)

const fileNameDB = "dbHandler.go"

// DbHandler holds the context and routes for this handler.
type DbHandler struct {
	Ctx         Context
	InsertTrack string
	GetTracks   string
	GetTrack    string
	DeleteTrack string
}

// Bind sets up the routes to the mux router.
func (dbHandler DbHandler) Bind(r *mux.Router) {
	r.HandleFunc(dbHandler.InsertTrack, dbHandler.insertTrackRecordPage).Methods(constant.Post)
	r.HandleFunc(dbHandler.GetTracks, dbHandler.getTracksPage).Methods(constant.Get)
	r.HandleFunc(dbHandler.GetTrack, dbHandler.getTrackPage).Queries("trID", "{trID}")
	r.HandleFunc(dbHandler.DeleteTrack, dbHandler.deleteTrackPage).Queries("trID", "{trID}")
}

// insertTrackRecordPage takes care of the overall logic of getting the request file saved
// and inserted into the DB.
func (dbHandler DbHandler) insertTrackRecordPage(w http.ResponseWriter, r *http.Request) {
	c, n, err := ProcessUploadRequest(r)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "insertTrackRecordPage", err)

		http.Error(w, err.Error(), c)
	}

	isPrivate := r.FormValue("private")
	bp := gtbackend.GetBoolFromString(isPrivate)

	err = insertTrackRecord(dbHandler.Ctx.App, n, bp)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "insertTrackRecordPage", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(c)
}

// getTracksPage retrieves a page of track metadata for the user
func (dbHandler DbHandler) getTracksPage(w http.ResponseWriter, r *http.Request) {
	// Extract data from header
	uId := r.Header.Get("UserID")
	pg := r.Header.Get("Page")
	qt := r.Header.Get("QueryType")
	ordDir := r.Header.Get("OrderDirection")

	// Process request
	d, err := GetTracks(dbHandler.Ctx.App, models.NewFirebaseQuery(uId, pg, qt, "Time", ordDir))
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "getTracksPage", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Convert to JSON
	rd, err := json.Marshal(d)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "getTracksPage", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(rd)
	return
}

func (dbHandler DbHandler) getTrackPage(w http.ResponseWriter, r *http.Request) {

}

func (dbHandler DbHandler) deleteTrackPage(w http.ResponseWriter, r *http.Request) {

}

// insertTrackRecord saves a FilePayload struct to the DB.
func insertTrackRecord(app *firebase.App, record models.FilePayload, isPrivate bool) (err error) {
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

	pIGC := parser.Parse()

	md := &models.IgcMetadata{
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

func getTrack(app *firebase.App, trackID string) (err error){
	return
}

func deleteTrack(app *firebase.App, trackID string) (err error) {
	return
}
