package models

import (
	"cloud.google.com/go/firestore"
	"github.com/GlidingTracks/gt-backend"
	"golang.org/x/tools/container/intsets"
	"strconv"
)

const filenameFQ = "FirebaseQuery.go"

// FirebaseQuery - struct denoting a Firebase query to receive data.
// UID string User ID as string
// TmSk int Time to skip the query to (pagination)
// Qt string Query type, if multiple flavors of queries exist (ex. getTracks personal and public)
// OrdDir firestore.Direction Which direction to order by. Default: Asc.
type FirebaseQuery struct {
	UID    string
	TmSk   int
	Qt     string
	OrdDir firestore.Direction
}

// NewFirebaseQuery - Initializes the query with values from strings (from header),
// sets TmSk = 1 and OrdDir = Asc as default.
func NewFirebaseQuery(u string, t string, q string, od string) FirebaseQuery {
	odfd := firestore.Asc
	if od == "Desc" {
		odfd = firestore.Desc
	}

	skipint := -1;
	var err error

	if t != "" {
		skipint, err = strconv.Atoi(t)
	}

	if err != nil || skipint < 1 {
		// Ensure that failed default timeskip will show the first page
		// Mismatch will lead to no results as all are skipped by Tmsk
		// Desc starts with largest time, Asc with smallest time
		if odfd == firestore.Desc {
			skipint = intsets.MaxInt
		} else {
			skipint = intsets.MinInt
		}

		gtbackend.DebugLog(gtbackend.InternalLog{
			Origin: filenameFQ,
			Method: "getTracksPage",
			Err:    err,
		})
	}

	return FirebaseQuery{u, skipint, q, odfd}
}
