package rest

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"google.golang.org/api/iterator"
)

// GetTracks gets a list of IgcMetadata from Firebase based on query
func GetTracks(app *firebase.App, query models.FirebaseQuery) (data []models.IgcMetadata, err error){
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "GetTracks", err)
		return
	}

	// Start query to Firebase based on the Query Type
	if query.Qt == "Private" {
		iter := client.Collection(constant.IgcMetadata).
			Where("UID", "==", query.UID).
			OrderBy(query.Ord, query.OrdDir).Documents(ctx)
		return processIterGetTracks(iter, query.Pg, "")
	} else {
		iter := client.Collection(constant.IgcMetadata).
			Where("Privacy", "==", false).
			OrderBy(query.Ord, query.OrdDir).Documents(ctx)
		return processIterGetTracks(iter, query.Pg, query.UID)
	}

	return data, err
}

/** processIterGetTracks
	Processes the request made to Firebase based
	iter *firestore.DocumentIterator Iterator with the results from firestore
	pg int Page to retrieve
	filterUID string Filter UID to remove from the results
 */
func processIterGetTracks(
	iter *firestore.DocumentIterator,
	pg int,
	filterUID string) (
	data []models.IgcMetadata,
	err error) {
	pageItemSkip := (pg - 1) * constant.PageSize

	// Process track query until length of data is the size of a page
	for ; len(data) < constant.PageSize; {
		doc, err := iter.Next()

		// Early break if there is no more data (last page)
		if err == iterator.Done {
			break
		}
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "processIterGetTracks", err)
			return data, err
		}

		// Convert doc to our model
		var d models.IgcMetadata
		err = doc.DataTo(&d)
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "processIterGetTracks", err)
			return data, err
		}

		// Filter out matching UID and add to data
		if d.UID != filterUID {
			if pageItemSkip > 0 {
				pageItemSkip--
			} else {
				data = append(data, d)
			}
		}
	}

	return data, err
}
