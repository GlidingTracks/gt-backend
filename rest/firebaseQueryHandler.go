package rest

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"net/http"
)

// fileNameFQH filenameDBH
const fileNameFQH = "firebaseQueryHandler.go"

// GetTracks gets a list of IgcMetadata from Firebase based on query
func GetTracks(app *firebase.App, query models.FirebaseQuery) (data []models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "GetTracks")
	ctx := context.Background()

	if app == nil {
		err = errors.New("Could not contact DB")
		return
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	// Start query to Firebase based on the Query Type
	if query.Qt == "Private" {
		iter := client.Collection(constant.IgcMetadata).
			Where("UID", "==", query.UID).
			OrderBy(constant.FirebaseQueryOrder, query.OrdDir).
			StartAfter(query.TmSk).Documents(ctx)
		return processIterGetTracks(iter, "")
	} else {
		iter := client.Collection(constant.IgcMetadata).
			Where("Privacy", "==", false).
			OrderBy(constant.FirebaseQueryOrder, query.OrdDir).
			StartAfter(query.TmSk).Documents(ctx)
		return processIterGetTracks(iter, query.UID)
	}

	return data, err
}

// GetTrack gets a track file from the Firebase Storage based on TrackID in metadata
func GetTrack(app *firebase.App, trackID string, uid string) (data []byte, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "GetTrack")

	// Verify that the user can download the file and abort if not
	_, d, err := getTrackMetadata(app, trackID, uid, false)
	if d.Privacy == true && d.Uid != uid {
		err = errors.New(constant.ErrorForbidden)
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

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

	// Read the entire file to data
	rc, err := bucket.Object(trackID).NewReader(context.Background())
	data, err = ioutil.ReadAll(rc)
	defer rc.Close()

	return
}

// DeleteTrack deletes the track from storage and firestore
func DeleteTrack(app *firebase.App, trackID string, uid string) (httpCode int, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "DeleteTrack")
	ctx := context.Background()
	httpCode = http.StatusBadRequest // Return before OK means failure

	// Verify that the user actually can delete this track
	client, _, err := getTrackMetadata(app, trackID, uid, true)
	if err != nil {
		httpCode = http.StatusForbidden
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	// Delete file from storage
	storageClient, err := app.Storage(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	bucket, err := storageClient.DefaultBucket()
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	err = bucket.Object(trackID).Delete(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Delete(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return
	}

	httpCode = http.StatusOK
	return
}

// UpdatePrivacy Updates privacy setting to new variable
func UpdatePrivacy(app *firebase.App, trackID string, uid string, newSetting bool) (updated models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "UpdatePrivacy")

	// Can only update privacy on owned tracks
	client, updated, err := getTrackMetadata(app, trackID, uid, true)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	// Set privacy and update it on firestore
	updated.Privacy = newSetting
	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Set(context.Background(), updated)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
	}

	return
}

// TakeOwnership Takes ownership of a track owned by the ScraperUID
func TakeOwnership(app *firebase.App, trackID string, uid string) (own models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "TakeOwnership")

	// Verify that the Scraper owns the track
	client, own, err := getTrackMetadata(app, trackID, constant.ScraperUID, true)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	// Take ownership and update status on firestore
	own.Uid = uid
	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Set(context.Background(), own)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
	}

	return
}

// InsertTrackPoint Inserts TrackPoint data for caching in the database
func InsertTrackPoint(app *firebase.App, trackID string, uid string, data []models.TrackPoint) (updated models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "InsertTrackPoint")
	client, updated, err := getTrackMetadata(app, trackID, uid, false)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	updated.TrackPoints = data

	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Set(context.Background(), updated)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
	}

	return
}

// Gets a single track metadata from firestore, can optionally verify that UID matches for security
func getTrackMetadata(app *firebase.App, trackID string, uid string, verifyUID bool) (client *firestore.Client, data models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "getTrackMetadata")
	ctx := context.Background()

	client, err = app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
	}

	doc, err := client.Collection(constant.IgcMetadata).Doc(trackID).Get(ctx)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	d, err := documentToModel(doc)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	if verifyUID && d.Uid != uid {
		err = errors.New(constant.ErrorForbidden)
		gtbackend.DebugLogErrNoMsg(log, err)
		return
	}

	data = d

	return
}

// documentToModel Pulls an IgcMetadata object out of a DocumentSnapshot
func documentToModel(doc *firestore.DocumentSnapshot) (data models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "documentToModel")
	// Convert doc to our model
	err = doc.DataTo(&data)
	if err != nil {
		gtbackend.DebugLogErrNoMsg(log, err)

		return data, err
	}
	return
}

/** processIterGetTracks
Processes the request made to Firebase based
iter *firestore.DocumentIterator Iterator with the results from firestore
filterUID string Filter UID to remove from the results
*/
func processIterGetTracks(iter *firestore.DocumentIterator, filterUID string) (data []models.IgcMetadata, err error) {
	log := gtbackend.DebugLogPrepareHeader(fileNameFQH, "processIterGetTracks")
	// Process track query until length of data is the size of a page
	for len(data) < constant.PageSize {
		doc, err := iter.Next()

		// Early break if there is no more data (last page)
		if err == iterator.Done {
			break
		}
		if err != nil {
			gtbackend.DebugLogErrNoMsg(log, err)

			return data, err
		}

		d, err := documentToModel(doc)
		if err != nil {
			gtbackend.DebugLogErrNoMsg(log, err)

			return data, err
		}

		// Filter out matching UID and add to data
		if d.Uid != filterUID && d.Uid != "" {
			data = append(data, d)
		}
	}

	return data, err
}
