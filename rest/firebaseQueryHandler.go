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

// fileNameFQH filename
const fileNameFQH = "firebaseQueryHandler.go"

// GetTracks gets a list of IgcMetadata from Firebase based on query
func GetTracks(app *firebase.App, query models.FirebaseQuery) (data []models.IgcMetadata, err error) {
	ctx := context.Background()

	if app == nil {
		err = errors.New("Could not contact DB")
		return
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "getTracks", Err: err})

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
func GetTrack(app *firebase.App, trackID string) (data []byte, err error) {
	client, err := app.Storage(context.Background())
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "GetTrack", Err: err})

		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "GetTrack", Err: err})

		return
	}

	// Read the entire file to data
	rc, err := bucket.Object(trackID).NewReader(context.Background())
	data, err = ioutil.ReadAll(rc)
	defer rc.Close()

	return
}

// DeleteTrack deletes the track from storage and firestore
func DeleteTrack(app *firebase.App, trackID string) (httpCode int, err error) {
	ctx := context.Background()
	httpCode = http.StatusBadRequest // Return before OK means failure

	// Delete file from storage
	storageClient, err := app.Storage(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "DeleteTrack", Err: err, Msg: "StorageClient"})

		return
	}

	bucket, err := storageClient.DefaultBucket()
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "DeleteTrack", Err: err, Msg: "Bucket"})

		return
	}

	err = bucket.Object(trackID).Delete(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "DeleteTrack", Err: err, Msg: "FileDelete"})

		return
	}

	// Delete file from firestore
	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "DeleteTrack", Err: err, Msg: "FirestoreClient"})

		return
	}

	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Delete(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "DeleteTrack", Err: err, Msg: "MetadataDelete"})

		return
	}

	httpCode = http.StatusOK
	return
}

func UpdatePrivacy(app *firebase.App, trackID string, uid string, newSetting bool) (updated models.IgcMetadata, err error) {
	client, ctx, err := getFirebaseClient(app)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "UpdatePrivacy", Err: err, Msg: "FirestoreClient"})
		return
	}

	updated, err = getTrackMetadata(client, trackID, uid, ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "UpdatePrivacy", Err: err, Msg: "Fail get metadata"})
		return
	}

	updated.Privacy = newSetting
	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Set(ctx, updated)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "UpdatePrivacy", Err: err, Msg: "Fail set metadata"})
	}

	return
}

func TakeOwnership(app *firebase.App, trackID string, uid string) (own models.IgcMetadata, err error) {
	client, ctx, err := getFirebaseClient(app)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "TakeOwnership", Err: err, Msg: "FirestoreClient"})
		return
	}

	own, err = getTrackMetadata(client, trackID, constant.ScraperUID, ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "TakeOwnership", Err: err, Msg: "Fail get metadata"})
		return
	}

	own.UID = uid
	_, err = client.Collection(constant.IgcMetadata).Doc(trackID).Set(ctx, own)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "TakeOwnership", Err: err, Msg: "Fail set metadata"})
	}

	return
}

func getFirebaseClient(app *firebase.App) (client *firestore.Client, ctx context.Context, err error) {
	ctx = context.Background()

	client, err = app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "getFirebaseClient", Err: err, Msg: "FirestoreClient"})
	}
	return
}

func getTrackMetadata(client *firestore.Client, trackID string, uid string, ctx context.Context) (data models.IgcMetadata, err error) {
	doc, err := client.Collection(constant.IgcMetadata).Doc(trackID).Get(ctx)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "getTrackMetadata", Err: err, Msg: "FirestoreClient"})
		return
	}

	d, err := documentToModel(doc)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "getTrackMetadata", Err: err, Msg: "FirestoreClient"})
		return
	}

	if d.UID != uid {
		err = errors.New(constant.ErrorForbidden)
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "getTrackMetadata", Err: err, Msg: "FirestoreClient"})
		return
	}

	data = d

	return
}

func documentToModel(doc *firestore.DocumentSnapshot) (data models.IgcMetadata, err error) {
	// Convert doc to our model
	err = doc.DataTo(&data)
	if err != nil {
		gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "processIterGetTracks", Err: err})

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
	// Process track query until length of data is the size of a page
	for len(data) < constant.PageSize {
		doc, err := iter.Next()

		// Early break if there is no more data (last page)
		if err == iterator.Done {
			break
		}
		if err != nil {
			gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "processIterGetTracks", Err: err})

			return data, err
		}

		d, err := documentToModel(doc)
		if err != nil {
			gtbackend.DebugLog(gtbackend.InternalLog{Origin: fileNameFQH, Method: "processIterGetTracks", Err: err})

			return data, err
		}

		// Filter out matching UID and add to data
		if d.UID != filterUID && d.UID != "" {
			data = append(data, d)
		}
	}

	return data, err
}
