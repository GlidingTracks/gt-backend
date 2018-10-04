package models

import (
	"cloud.google.com/go/firestore"
	"github.com/GlidingTracks/gt-backend"
	"strconv"
)

const filenameFQ = "FirebaseQuery.go"

/**
Struct denoting a Firebase query to receive data.
UID string User ID as string
Pg int Page number to get
Qt string Query type, if multiple flavors of queries exist (ex. getTracks personal and public)
Ord string Which data type to order by
OrdDir firestore.Direction Which direction to order by. Default: Asc
*/
type FirebaseQuery struct {
	UID    string
	Pg     int
	Qt     string
	Ord    string
	OrdDir firestore.Direction
}

func NewFirebaseQuery(u string, p string, q string, o string, od string) FirebaseQuery {
	pint, err := strconv.Atoi(p)
	if err != nil {
		pint = 1

		gtbackend.DebugLog(filenameFQ, "getTracksPage", err)
	}

	odfd := firestore.Asc
	if od == "Desc" {
		odfd = firestore.Desc
	}

	return FirebaseQuery{u, pint, q, o, odfd}
}
