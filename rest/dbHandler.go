package rest

import (
	"encoding/json"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

const fileNameDB = "dbHandler.go"

// DbHandler holds the context and routes for this handler.
type DbHandler struct {
	Ctx           Context
	InsertTrack   string
	GetTracks     string
	GetTrack      string
	DeleteTrack   string
	UpdatePrivacy string
}

// Bind sets up the routes to the mux router.
func (dbHandler DbHandler) Bind(r *mux.Router) {
	r.HandleFunc(dbHandler.InsertTrack, dbHandler.insertTrackRecordPage).Methods(constant.Post)
	r.HandleFunc(dbHandler.GetTracks, dbHandler.getTracksPage).Methods(constant.Get)
	r.HandleFunc(dbHandler.GetTrack, dbHandler.getTrackPage).Methods(constant.Get)
	r.HandleFunc(dbHandler.DeleteTrack, dbHandler.deleteTrackPage).Methods(constant.Delete)
	r.HandleFunc(dbHandler.UpdatePrivacy, dbHandler.updatePrivacyPage).Methods(constant.Put)
}

// insertTrackRecordPage takes care of the overall logic of getting the request file saved
// and inserted into the DB.
func (dbHandler DbHandler) insertTrackRecordPage(w http.ResponseWriter, r *http.Request) {
	c, md, err := ProcessUploadRequest(dbHandler.Ctx.App, r)
	if err != nil {
		http.Error(w, err.Error(), c)
		return
	}

	// Convert to JSON and prepare response
	d, err := json.Marshal(md)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "insertTrack",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w = prepareGeneralResponse(w, d, constant.ApplicationJSON)
	return
}

// getTracksPage retrieves a page of track metadata for the user
func (dbHandler DbHandler) getTracksPage(w http.ResponseWriter, r *http.Request) {
	// Extract data from header
	uID := r.Header.Get("uid")
	tmsk := r.Header.Get("timeSkip")
	qt := r.Header.Get("queryType")
	ordDir := r.Header.Get("orderDirection")

	if uID == "" {
		http.Error(w, errors.New("No uId supplied").Error(), http.StatusBadRequest)
	}

	// Process request
	d, err := GetTracks(dbHandler.Ctx.App, models.NewFirebaseQuery(uID, tmsk, qt, ordDir))
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "getTracksPage",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to JSON and prepare response
	rd, err := json.Marshal(d)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "getTracks",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w = prepareGeneralResponse(w, rd, constant.ApplicationJSON)
	return
}

// getTrackPage Gets the track from the Firebase Storage based on TrackID
func (dbHandler DbHandler) getTrackPage(w http.ResponseWriter, r *http.Request) {
	trackID := r.Header.Get("trackID")

	// Process request
	d, err := GetTrack(dbHandler.Ctx.App, trackID)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "getTrack",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	w = prepareGeneralResponse(w, d, constant.TextPlain)
	return
}

// deleteTrackPage Deletes track from Storage and Firestore Database based on TrackID
func (dbHandler DbHandler) deleteTrackPage(w http.ResponseWriter, r *http.Request) {
	trackID := r.Header.Get("trackID")

	// Process request
	c, err := DeleteTrack(dbHandler.Ctx.App, trackID)
	if err != nil || c != http.StatusOK {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "deleteTrackPage",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w = prepareGeneralResponse(w, []byte(trackID), constant.TextPlain)
	return
}

func (dbHandler DbHandler) updatePrivacyPage(w http.ResponseWriter, r *http.Request) {
	trackID := r.Header.Get("trackID")
	newSetting := gtbackend.GetBoolFromString(r.Header.Get("private"))
	uid := r.Header.Get("uid")

	d, err := UpdatePrivacy(dbHandler.Ctx.App, trackID, uid, newSetting)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "updatePrivacyPage",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to JSON and prepare response
	rd, err := json.Marshal(d)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: fileNameDB,
			Method: "updatePrivacyPage",
			Err:    err,
		})

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w = prepareGeneralResponse(w, rd, constant.ApplicationJSON)
	return
}

// Prepares a general response to send back to the client, setting various common variables in the ResponseWriter
func prepareGeneralResponse(unPrepW http.ResponseWriter, rawData []byte, contentType string) (prepW http.ResponseWriter) {
	prepW = unPrepW

	prepW.Header().Set(constant.ContentType, contentType)
	prepW.WriteHeader(http.StatusOK)
	prepW.Write(rawData)
	return
}
