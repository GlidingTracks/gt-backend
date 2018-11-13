package rest

import (
	"encoding/json"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/gorilla/mux"
	"net/http"
)

const filenameDBH = "dbHandler.go"

// DbHandler holds the context and routes for this handler.
type DbHandler struct {
	Ctx              Context
	InsertTrack      string
	GetTracks        string
	GetTrack         string
	DeleteTrack      string
	UpdatePrivacy    string
	TakeOwnership    string
	InsertTrackPoint string
}

// Bind sets up the routes to the mux router.
func (dbHandler DbHandler) Bind(r *mux.Router) {
	r.HandleFunc(dbHandler.InsertTrack, dbHandler.insertTrackRecordPage).Methods(constant.Post)
	r.HandleFunc(dbHandler.GetTracks, dbHandler.getTracksPage).Methods(constant.Get)
	r.HandleFunc(dbHandler.GetTrack, dbHandler.getTrackPage).Methods(constant.Get)
	r.HandleFunc(dbHandler.DeleteTrack, dbHandler.deleteTrackPage).Methods(constant.Delete)
	r.HandleFunc(dbHandler.UpdatePrivacy, dbHandler.updatePrivacyPage).Methods(constant.Put)
	r.HandleFunc(dbHandler.TakeOwnership, dbHandler.takeOwnershipPage).Methods(constant.Put)
	r.HandleFunc(dbHandler.InsertTrackPoint, dbHandler.insertTrackPointPage).Methods(constant.Put)
}

// insertTrackRecordPage takes care of the overall logic of getting the request file saved
// and inserted into the DB.
func (dbHandler DbHandler) insertTrackRecordPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "insertTrackRecordPage")
	c, d, err := ProcessUploadRequest(dbHandler.Ctx.App, r)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		http.Error(w, err.Error(), c)
		return
	}

	w = prepareMetadataSingleResponse(w, d, constant.ApplicationJSON)
	return
}

// getTracksPage retrieves a page of track metadata for the user
func (dbHandler DbHandler) getTracksPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "getTracksPage")
	// Extract data from header
	uID := r.Header.Get("uid")
	tmsk := r.Header.Get("timeSkip")
	qt := r.Header.Get("queryType")
	ordDir := r.Header.Get("orderDirection")

	// Process request
	d, err := GetTracks(dbHandler.Ctx.App, models.NewFirebaseQuery(uID, tmsk, qt, ordDir))
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to JSON and prepare response
	rd, err := json.Marshal(d)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w = prepareGeneralResponse(w, rd, constant.ApplicationJSON)
	return
}

// getTrackPage Gets the track from the Firebase Storage based on TrackID
func (dbHandler DbHandler) getTrackPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "getTrackPage")
	trackID := r.Header.Get("trackID")
	uid := r.Header.Get("uid")

	// Process request
	d, err := GetTrack(dbHandler.Ctx.App, trackID, uid)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	w = prepareGeneralResponse(w, d, constant.TextPlain)
	return
}

// deleteTrackPage Deletes track from Storage and Firestore Database based on TrackID
func (dbHandler DbHandler) deleteTrackPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "deleteTrackPage")
	trackID := r.Header.Get("trackID")
	uid := r.Header.Get("uid")

	// Process request
	c, err := DeleteTrack(dbHandler.Ctx.App, trackID, uid)
	if err != nil || c != http.StatusOK {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w = prepareGeneralResponse(w, []byte(trackID), constant.TextPlain)
	return
}

// updatePrivacyPage Changes the Private variable to a new variable
func (dbHandler DbHandler) updatePrivacyPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "updatePrivacyPage")
	trackID := r.Header.Get("trackID")
	newSetting := gtbackend.GetBoolFromString(r.Header.Get("private"))
	uid := r.Header.Get("uid")

	d, err := UpdatePrivacy(dbHandler.Ctx.App, trackID, uid, newSetting)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w = prepareMetadataSingleResponse(w, d, constant.ApplicationJSON)
	return
}

// takeOwnershipPage Sets the UID of tracks to this UID if it's the scraper's UID that is taken from
func (dbHandler DbHandler) takeOwnershipPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "takeOwnershipPage")
	trackID := r.Header.Get("trackID")
	uid := r.Header.Get("uid")

	d, err := TakeOwnership(dbHandler.Ctx.App, trackID, uid)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w = prepareMetadataSingleResponse(w, d, constant.ApplicationJSON)
	return
}

// insertTrackPointPage Inserts TrackPoint data for caching in the database
func (dbHandler DbHandler) insertTrackPointPage(w http.ResponseWriter, r *http.Request) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "insertTrackPointPage")
	trackID := r.Header.Get("trackID")
	uid := r.Header.Get("uid")
	rawJSON := r.Header.Get("trackPoints")

	// Parse into objects the string JSON from header
	var parsed []models.TrackPoint
	err := json.Unmarshal([]byte(rawJSON), &parsed)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := InsertTrackPoint(dbHandler.Ctx.App, trackID, uid, parsed)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w = prepareMetadataSingleResponse(w, d, constant.ApplicationJSON)
	return
}

// prepareGeneralResponse Prepares a general response to send back to the client, setting various common variables in the ResponseWriter
func prepareGeneralResponse(unPrepW http.ResponseWriter, rawData []byte, contentType string) (prepW http.ResponseWriter) {
	prepW = unPrepW

	prepW.Header().Set(constant.ContentType, contentType)
	prepW.WriteHeader(http.StatusOK)
	prepW.Write(rawData)
	return
}

// prepareMetadataSingleResponse Prepares a metadata with a single IgcMetadata object response
func prepareMetadataSingleResponse(unPrepW http.ResponseWriter, data models.IgcMetadata, contentType string) (prepW http.ResponseWriter) {
	log := gtbackend.DebugLogPrepareHeader(filenameDBH, "prepareMetadataSingleResponse")
	prepW = unPrepW

	d, err := json.Marshal(data)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		http.Error(prepW, err.Error(), http.StatusBadRequest)
		return
	}

	prepW = prepareGeneralResponse(prepW, d, contentType)
	return
}
