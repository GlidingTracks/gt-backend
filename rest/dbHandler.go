package rest

import (
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"net/http"
)

const fileNameDB = "dbHandler.go"

// DbHandler holds the context and routes for this handler.
type DbHandler struct {
	Ctx         Context
	InsertTrack string
	GetTracks string
	GetTrack string
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

	err = insertTrackRecord(dbHandler.Ctx.App, n)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "insertTrackRecordPage", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(c)
}

func (dbHandler DbHandler) getTracksPage(w http.ResponseWriter, r *http.Request) {

	d, err := getTracks(dbHandler.Ctx.App)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "getTracksPage", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	rd, err := json.Marshal(d)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "getTracksPage", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}


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
func insertTrackRecord(app *firebase.App, record models.FilePayload) (err error) {
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		return
	}

	_, _, err = client.Collection(constant.CollectionTracks).Add(ctx, record)
	if err != nil {
		return
	}

	return
}

// getTracks gets a list of tracks from the DB
func getTracks(app *firebase.App) (d []models.IgcMetadata, err error){
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "getTracks", err)

		return
	}

	iter := client.Collection(constant.IgcMetadata).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "getTracks", err)
			return d, err
		}
		var c models.IgcMetadata

		doc.DataTo(&c)
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "getTracks", err)
			return d, err
		}

		d = append(d, c)

	}

	return
}

func getTrack(app *firebase.App, trackID string) (err error){
	return
}

func deleteTrack(app *firebase.App, trackID string) (err error){
	return
}
