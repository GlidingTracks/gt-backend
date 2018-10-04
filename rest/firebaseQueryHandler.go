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

// getTracks gets a list of tracks from the DB
func GetTracks(app *firebase.App, query models.FirebaseQuery) (data []models.IgcMetadata, err error) {
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		gtbackend.DebugLog(fileNameDB, "getTracks", err)

		return
	}

	println(query.Qt + " - " + query.UID)

	if query.Qt == "Personal" {
		iter := client.Collection(constant.IgcMetadata).
			Where("UID", "==", query.UID).
			OrderBy(query.Ord, query.OrdDir).Documents(ctx)
		return processIterGetTracks(iter, query, false)
	} else {
		iter := client.Collection(constant.IgcMetadata).
			Where("Privacy", "==", false).
			OrderBy(query.Ord, query.OrdDir).Documents(ctx)
		return processIterGetTracks(iter, query, true)
	}

	return data, err
}

func processIterGetTracks(iter *firestore.DocumentIterator, query models.FirebaseQuery, filterSelf bool) (data []models.IgcMetadata, err error) {
	total := 0
	for ; total < (query.Pg-1)*constant.PageSize; total++ {
		iter.Next()
		if err == iterator.Done {
			return data, err
		}
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "getTracks", err)
			return data, err
		}
	}

	for ; total < query.Pg*constant.PageSize; total++ {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "getTracks", err)
			return data, err
		}

		var d models.IgcMetadata

		println("Object")

		doc.DataTo(&d)
		if err != nil {
			gtbackend.DebugLog(fileNameDB, "getTracks", err)
			return data, err
		}

		if filterSelf {
			if d.UID != query.UID {
				data = append(data, d)
			}
		} else {
			data = append(data, d)
		}

	}

	return data, err
}
