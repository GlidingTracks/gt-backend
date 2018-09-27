package rest

import (
	"context"
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
}

// Bind sets up the routes to the mux router.
func (dbHandler DbHandler) Bind(r *mux.Router) {
	r.HandleFunc("/insertTrack", dbHandler.insertTrackRecordPage).Methods(constant.Post)
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
