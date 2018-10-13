package rest

import (
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
	c, _, err := ProcessUploadRequest(dbHandler.Ctx.App, r)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "insertTrackRecordPage", err)

		http.Error(w, err.Error(), c)
		return
	}

	w.WriteHeader(c)
}

// getTracksPage retrieves a page of track metadata for the user
func (dbHandler DbHandler) getTracksPage(w http.ResponseWriter, r *http.Request) {
	// Extract data from header
	uID := r.Header.Get("uid")
	tmsk := r.Header.Get("timeSkip")
	qt := r.Header.Get("queryType")
	ordDir := r.Header.Get("orderDirection")

	// Process request
	d, err := GetTracks(dbHandler.Ctx.App, models.NewFirebaseQuery(uID, tmsk, qt, ordDir))
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

func getTrack(app *firebase.App, trackID string) (err error) {
	return
}

func deleteTrack(app *firebase.App, trackID string) (err error) {
	return
}
